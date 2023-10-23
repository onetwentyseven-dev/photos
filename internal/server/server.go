package server

import (
	"context"
	"net/http"
	"photos"
	"photos/internal"
	"photos/internal/server/router"
	"photos/internal/server/templates"
	"photos/internal/store/mysql"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ddouglas/authenticator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type Server struct {
	appURL string
	env    photos.Environment
	port   string

	http   *http.Server
	router *router.Service

	authenticator *authenticator.Service
	decoder       *schema.Decoder
	logger        *logrus.Logger
	s3            *s3.Client
	sessions      sessions.Store
	templates     *templates.Service

	userRepository  *mysql.UserRepository
	imageRepository *mysql.ImageRepository
}

func New(
	appURL string,
	env photos.Environment,
	port string,
	logger *logrus.Logger,
	authenticator *authenticator.Service,
	s3Client *s3.Client,
	sessionStore sessions.Store,
	userRepository *mysql.UserRepository,
	imageRepository *mysql.ImageRepository,
) *Server {
	s := &Server{
		appURL: appURL,
		env:    env,
		port:   port,

		authenticator: authenticator,
		decoder:       schema.NewDecoder(),
		logger:        logger,
		s3:            s3Client,
		sessions:      sessionStore,

		userRepository:  userRepository,
		imageRepository: imageRepository,
	}

	s.router = s.buildRouter()
	s.templates = templates.New(s.router, logger)
	s.http = &http.Server{
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		Addr:         ":" + port,
		Handler:      s.router.PrimaryRouter(),
	}

	return s

}
func (s *Server) Start() error {
	s.logger.WithField("port", s.port).Info("starting server")
	return s.http.ListenAndServe()
}

func (s *Server) Mux() *mux.Router {
	return s.router.PrimaryRouter()
}

func (s *Server) GracefullyShutdown(ctx context.Context) error {
	s.logger.Info("stopping server")
	return s.http.Shutdown(ctx)
}

func (s *Server) buildRouter() *router.Service {

	r := router.New()
	r.Use(router.PrimaryRouter, s.logging, s.user)
	r.PrimaryRouter().PathPrefix("/static").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("cache-control", "max-age=86400")
		http.StripPrefix("/static/", http.FileServer(http.FS(photos.AssetFS(s.env)))).ServeHTTP(w, r)
	})).Name("static").Methods(http.MethodGet)

	r.AddRoute(router.PrimaryRouter, "home", http.MethodGet, "/", s.handleHome)
	r.AddRoute(router.PrimaryRouter, "login", http.MethodGet, "/login", s.handleLogin)
	r.AddRoute(router.PrimaryRouter, "logout", http.MethodGet, "/logout", s.handleLogout)
	r.AddRoute(router.PrimaryRouter, "oauth-callback", http.MethodGet, "/oauth/callback", s.handleOauthCallback)

	r.NewSubRouter(router.APIRouter)
	r.Use(router.APIRouter, s.user, s.apiAuth)
	r.AddRoute(router.APIRouter, "api-auth-validate", http.MethodGet, "/api/auth/validate", s.handleValidateAuth)
	r.AddRoute(router.APIRouter, "api-image-metadata", http.MethodPost, "/api/image/metadata", s.handlePostImageMeta)

	r.NewSubRouter(router.AuthRouter)
	r.Use(router.AuthRouter, s.auth)

	r.AddRoute(router.AuthRouter, "dashboard", http.MethodGet, "/dashboard", s.handleDashboard)
	r.AddRoute(router.AuthRouter, "upload", http.MethodGet, "/dashboard/upload", s.handleGetUpload)

	return r
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()
	user := internal.UserFromContext(ctx)

	err := s.templates.Homepage(ctx, user).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render homepage")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *Server) writeRedirectRouteName(w http.ResponseWriter, routeName string, routePairs ...string) {

	entry := s.logger.WithField("routeName", routeName)

	route := s.router.PrimaryRouter().GetRoute(routeName)
	if route == nil {
		entry.Error("no route found for name, unable to perform redirect")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uri, err := route.URL(routePairs...)
	if err != nil {
		entry.WithError(err).Error("failed to generate uri")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", uri.String())
	w.WriteHeader(http.StatusTemporaryRedirect)

}

package server

import (
	"fmt"
	"net/http"
	"photos/internal"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (s *Server) logging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		entry := logrus.NewEntry(s.logger)
		method, path := r.Method, r.URL.Path
		var wrapped = &wrappedResponseWriter{
			ResponseWriter: w,
		}
		handler.ServeHTTP(wrapped, r)
		var status = wrapped.status
		if status == 0 {
			status = http.StatusOK
		}

		entry.WithFields(logrus.Fields{
			"duration": time.Since(start),
			"status":   status,
		}).Infof(fmt.Sprintf("%s %s", method, path))

	})
}

func (s *Server) user(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ctx = r.Context()

		session, err := s.sesssionFromRequest(r)
		if err != nil {
			// Create an error page and redirect to that. Use session flashing to flash an internal error message of sorts
			s.logger.WithError(err).Error("failed to load session")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userIDInf, ok := session.Values["userID"]
		if !ok {
			handler.ServeHTTP(w, r)
			return
		}

		userID, err := uuid.Parse(userIDInf.(string))
		if err != nil {
			s.logger.WithError(err).Error("failed to parse userID to valid uuid")
			// s.writeRedirectRouteName(w, "login")
			handler.ServeHTTP(w, r)
			return
		}

		user, err := s.userRepository.User(ctx, userID)
		if err != nil {
			s.logger.WithError(err).Error("failed to look up user by id")
			// s.writeRedirectRouteName(w, "login")
			handler.ServeHTTP(w, r)
			return
		}

		ctx = internal.ContextWithUser(ctx, user)

		handler.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (s *Server) auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ctx = r.Context()

		user := internal.UserFromContext(ctx)
		if user == nil {
			s.logger.Error("no user found in context, redirecting")
			s.writeRedirectRouteName(w, "login")
			return
		}

		handler.ServeHTTP(w, r)

	})
}

func (s *Server) apiAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ctx = r.Context()

		user := internal.UserFromContext(ctx)
		if user == nil {
			s.logger.Error("no user found in context, redirecting")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)

	})
}

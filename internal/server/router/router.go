package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type RouterKey string

const (
	PrimaryRouter RouterKey = "primary"
	AuthRouter    RouterKey = "auth"
	APIRouter     RouterKey = "api"
)

type Service struct {
	routers map[RouterKey]*mux.Router
}

func New() *Service {
	return &Service{
		routers: map[RouterKey]*mux.Router{
			PrimaryRouter: mux.NewRouter(),
		},
	}
}

// func AddRoute(router *mux.Router, name, method, path string, f http.HandlerFunc) {
// 	router.HandleFunc(path, f).Methods(method).Name(name)
// }

func (s *Service) NewSubRouter(key RouterKey) {
	_, ok := s.routers[key]
	if ok {
		panic(fmt.Sprintf("router key already exists: %s", key))
	}
	s.routers[key] = s.routers[PrimaryRouter].NewRoute().Subrouter()
}

func (s *Service) AddRoute(router RouterKey, name, method, path string, f http.HandlerFunc) {
	s.routers[router].HandleFunc(path, f).Methods(method).Name(name)
}

func (s *Service) Use(router RouterKey, mwf ...mux.MiddlewareFunc) {

	_, ok := s.routers[router]
	if !ok {
		panic(fmt.Sprintf("unknown router key: %s", router))
	}

	s.routers[router].Use(mwf...)
}

func (s *Service) GetRouter(key RouterKey) *mux.Router {
	return s.routers[key]
}

func (s *Service) PrimaryRouter() *mux.Router {
	return s.routers[PrimaryRouter]
}

func (s *Service) BuildURI(name string, pairs ...string) string {
	route := s.PrimaryRouter().Get(name)
	if route == nil {
		return ""
	}

	uri, err := route.URL(pairs...)
	if err != nil {
		return ""
	}

	return uri.String()
}

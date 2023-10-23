package templates

import (
	"photos/internal/server/router"

	"github.com/sirupsen/logrus"
)

type Service struct {
	router *router.Service
	logger *logrus.Logger
}

func New(
	router *router.Service,
	logger *logrus.Logger,
) *Service {
	return &Service{
		logger: logger,
		router: router,
	}
}

func (s *Service) buildRoute(name string, pairs ...string) string {
	return s.router.BuildURI(name, pairs...)
}

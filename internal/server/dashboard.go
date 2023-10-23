package server

import (
	"net/http"
	"photos/internal"
)

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var user = internal.UserFromContext(ctx)
	err := s.templates.Dashboard(ctx, user).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render dashboard")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

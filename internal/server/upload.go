package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"photos"
	"photos/internal"

	"github.com/google/uuid"
)

func (s *Server) handleGetUpload(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var user = internal.UserFromContext(ctx)

	err := s.templates.Upload(ctx, user).Render(w)
	if err != nil {
		s.logger.WithError(err).Error("failed to render upload")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

type imageMetaRequest struct {
	Name string `json:"name"`
}

func (s *Server) handlePostImageMeta(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var user = internal.UserFromContext(ctx)

	var payload = new(imageMetaRequest)

	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.WithError(err).Error("failed to read body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("REQUEST PAYLOAD :: ", string(data))

	err = json.Unmarshal(data, payload)
	if err != nil {
		s.logger.WithError(err).Error("failed to decode payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	image := &photos.Image{
		ID:     uuid.New().String(),
		UserID: user.ID,
		Name:   payload.Name,
		Status: photos.ProcessingImageStatus,
	}

	err = s.imageRepository.CreateImage(ctx, image)
	if err != nil {
		s.logger.WithError(err).Error("failed to create image")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(image)
	if err != nil {
		s.logger.WithError(err).Error("failed to encode image")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

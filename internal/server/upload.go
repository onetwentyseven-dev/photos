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

type postImageMetaRequest struct {
	Name string `json:"name"`
}

func (s *Server) handlePostImageMeta(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var user = internal.UserFromContext(ctx)

	var payload = new(postImageMetaRequest)

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
		Status: photos.QueuedImageStatus,
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

type patchImageMetaRequest struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Error  *string `json:"error"`
}

func (s *Server) handlePatchImageMeta(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var user = internal.UserFromContext(ctx)

	var payload = new(patchImageMetaRequest)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		s.logger.WithError(err).Error("failed to decode payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	image, err := s.imageRepository.Image(ctx, payload.ID)
	if err != nil {
		s.logger.WithError(err).Error("failed to get image")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if image == nil {
		s.logger.WithError(err).Error("image not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if image.UserID != user.ID {
		s.logger.WithError(err).Error("image not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	status := photos.ImageStatus(payload.Status)
	if !status.Valid() {
		s.logger.WithError(err).Error("invalid image status")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	image.Status = status
	if status == photos.ErroredImageStatus {
		image.ProcessingErrors = append(image.ProcessingErrors, *payload.Error)
	}

	err = s.imageRepository.UpdateImage(ctx, image)
	if err != nil {
		s.logger.WithError(err).Error("failed to update image")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func l(m string, i ...any) {

	data, _ := json.Marshal(i)
	fmt.Printf("%s :: %s", m, string(data))

}

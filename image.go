package photos

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type ImageStatus string

const (
	ErroredImageStatus    ImageStatus = "errored"
	QueuedImageStatus     ImageStatus = "queued"
	ProcessingImageStatus ImageStatus = "processing"
	ProcessedImageStatus  ImageStatus = "processed"
)

var AllImageStatuses = []ImageStatus{
	ErroredImageStatus, QueuedImageStatus,
	ProcessingImageStatus, ProcessedImageStatus,
}

func (is ImageStatus) Valid() bool {
	for _, status := range AllImageStatuses {
		if status == is {
			return true
		}
	}
	return false
}

type ImageProcessingErrors []string

func (ipe ImageProcessingErrors) Scan(value interface{}) error {

	if value == nil {
		return nil
	}

	return json.Unmarshal(value.([]byte), &ipe)

}

func (ipe ImageProcessingErrors) Value() (driver.Value, error) {
	if len(ipe) == 0 {
		return []byte(`[]`), nil
	}

	return json.Marshal(ipe)
}

type ImageExifData map[string]any

func (ied ImageExifData) Scan(value interface{}) error {

	if value == nil {
		return nil
	}

	return json.Unmarshal(value.([]byte), &ied)

}

func (ied ImageExifData) Value() (driver.Value, error) {
	if len(ied) == 0 {
		return []byte(`{}`), nil
	}

	return json.Marshal(ied)
}

type Image struct {
	ID               string                `structs:"id" db:"id" json:"id"`
	UserID           string                `structs:"user_id" db:"user_id" json:"user_id"`
	Name             string                `structs:"name" db:"name" json:"name"`
	Description      *string               `structs:"description" db:"description" json:"description,omitempty"`
	Status           ImageStatus           `structs:"status" db:"status" json:"status"`
	ProcessingErrors ImageProcessingErrors `structs:"processing_errors" db:"processing_errors" json:"processing_errors,omitempty"`
	ExifData         ImageExifData         `structs:"image_exif_data" db:"image_exif_data" json:"image_exif_data,omitempty"`
	TSCreated        time.Time             `structs:"ts_created" db:"ts_created" json:"ts_created,omitempty"`
	TSUpdated        time.Time             `structs:"ts_updated" db:"ts_updated" json:"ts_updated,omitempty"`
}

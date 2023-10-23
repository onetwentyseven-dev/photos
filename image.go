package photos

import (
	"time"
)

type ImageStatus string

const (
	ProcessingImageStatus ImageStatus = "processing"
	ProcessedImageStatus  ImageStatus = "processed"
)

type Image struct {
	ID          string      `structs:"id" db:"id" json:"id"`
	UserID      string      `structs:"user_id" db:"user_id" json:"user_id"`
	Name        string      `structs:"name" db:"name" json:"name"`
	Description *string     `structs:"description" db:"description" json:"description,omitempty"`
	Status      ImageStatus `structs:"status" db:"status" json:"status"`
	TSCreated   time.Time   `structs:"ts_created" db:"ts_created" json:"ts_created,omitempty"`
	TSUpdated   time.Time   `structs:"ts_updated" db:"ts_updated" json:"ts_updated,omitempty"`
}

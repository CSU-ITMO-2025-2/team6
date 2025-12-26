package model

import (
	"time"

	"github.com/google/uuid"
)

type Study struct {
	ID               uuid.UUID  `json:"id"`
	Name             *string    `json:"name"`
	Status           string     `json:"status"`
	Image            *Image     `json:"image"`
	OwnerID          uuid.UUID  `json:"owner_id"`
	ImageID          uuid.UUID  `json:"image_id"`
	ImageURL         *string    `json:"image_url"`
	PredictedClassID *uuid.UUID `json:"predicted_class_id"`
	PredictedClass   *string    `json:"predicted_class"`
	PredictedScore   *float64   `json:"predicted_score"`
	ErrorDescription *string    `json:"error_description"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

type Image struct {
	ChunkData []byte `json:"chunk_data"`
	FileName  string `json:"file_name"`
	MimeType  string `json:"mime_type"`
}

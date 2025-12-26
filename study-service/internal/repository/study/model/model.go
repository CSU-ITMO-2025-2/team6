package model

import (
	"time"

	"github.com/google/uuid"
)

type Study struct {
	ID               uuid.UUID  `json:"id"`
	Name             *string    `json:"name"`
	Status           string     `json:"status"`
	OwnerID          uuid.UUID  `json:"owner_id"`
	ImageID          uuid.UUID  `json:"image_id"`
	PredictedClassID *uuid.UUID `json:"predicted_class_id"`
	PredictedScore   *int       `json:"predicted_score"`
	ErrorDescription *string    `json:"error_description"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

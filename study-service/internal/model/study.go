package model

import (
	"time"

	"github.com/google/uuid"
)

type StudyStatusType string

const (
	New        StudyStatusType = "NEW"
	Queued     StudyStatusType = "QUEUED"
	Processing StudyStatusType = "PROCESSING"
	Completed  StudyStatusType = "COMPLETED"
	Failed     StudyStatusType = "FAILED"
)

type Study struct {
	ID               uuid.UUID       `json:"id"`
	Name             *string         `json:"name"`
	Status           StudyStatusType `json:"status"`
	OwnerID          uuid.UUID       `json:"owner_id"`
	ImageID          uuid.UUID       `json:"image_id"`
	ImageURL         *string         `json:"image_url"`
	PredictedClassID *uuid.UUID      `json:"predicted_class_id"`
	PredictedScore   *int            `json:"predicted_score"`
	ErrorDescription *string         `json:"error_description"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        *time.Time      `json:"updated_at"`
}

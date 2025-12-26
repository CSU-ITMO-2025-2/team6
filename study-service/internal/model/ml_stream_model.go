package model

import "github.com/google/uuid"

type MLStream struct {
	StudyID uuid.UUID `json:"study_id"`
	Image   []byte    `json:"image"`
}

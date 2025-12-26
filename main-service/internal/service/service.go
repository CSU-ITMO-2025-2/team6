package service

import (
	"context"

	"main-service/internal/model"

	"github.com/google/uuid"
)

type StudyService interface {
	Create(ctx context.Context, st *model.Study) (uuid.UUID, error)
}

type UserService interface {
}

package client

import (
	"context"

	"main-service/internal/model"

	"github.com/google/uuid"
)

type StudyClient interface {
	CreateStudy(ctx context.Context, st *model.Study) (uuid.UUID, error)
	GetStudy(ctx context.Context, id uuid.UUID) (*model.Study, error)
}

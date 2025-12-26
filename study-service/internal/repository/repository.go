package repository

import (
	"context"
	"io"

	"study-service/internal/model"

	"github.com/google/uuid"
)

type StudyRepository interface {
	Create(ctx context.Context, study *model.Study) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (*model.Study, error)
	Update(ctx context.Context, s *model.Study) (*model.Study, error)
	Delete(ctx context.Context, id uuid.UUID) error

	//List(ctx context.Context, filter *Filter, limit, offset int) ([]*model.Study, error)
}

type Storage interface {
	Upload(ctx context.Context, bucketName, objectName uuid.UUID, data []byte, contentType string) error
	Download(ctx context.Context, bucketName, objectName string) (io.Reader, error)
	Delete(ctx context.Context, bucketName, objectName string) error
	GetURL(bucketName, objectName string) string
}

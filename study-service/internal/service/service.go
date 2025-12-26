package service

import (
	"context"

	"local-lib/queue"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type StudyService interface {
	CreateStudy(ctx context.Context, userID uuid.UUID, image []byte, contentType string) (uuid.UUID, error)
}

type StudyAnalyzeWorkerService interface {
	ImageAnalysisWorker(log *zap.SugaredLogger, n queue.Study)
}

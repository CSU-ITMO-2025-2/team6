package workers

import (
	"context"

	"local-lib/queue"
	"study-service/internal/service"

	"go.uber.org/zap"
)

func StartWorkers(ctx context.Context, log *zap.SugaredLogger, q queue.Study, study service.StudyAnalyzeWorkerService) {
	go study.ImageAnalysisWorker(log, q)

	// Если появятся другие:
	// go study.SomeOtherWorker(log, q)
	// go study.CleanupWorker(log)
}

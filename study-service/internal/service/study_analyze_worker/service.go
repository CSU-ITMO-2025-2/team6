package study_analyze_worker

import (
	"local-lib/queue"
	"study-service/internal/repository"
	"study-service/internal/service"
)

type serv struct {
	studyRepository repository.StudyRepository
	storage         repository.Storage
	queue           queue.Study
}

func NewService(studyRepository repository.StudyRepository, storage repository.Storage, queue queue.Study) service.StudyAnalyzeWorkerService {
	return &serv{
		studyRepository: studyRepository,
		storage:         storage,
		queue:           queue,
	}
}

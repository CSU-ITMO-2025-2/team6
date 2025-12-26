package study

import (
	"local-lib/queue"
	"study-service/internal/service"

	"study-service/internal/repository"
)

type serv struct {
	studyRepository repository.StudyRepository
	storage         repository.Storage
	queue           queue.Study
}

func NewService(studyRepository repository.StudyRepository, storage repository.Storage, queue queue.Study) service.StudyService {
	return &serv{
		studyRepository: studyRepository,
		storage:         storage,
		queue:           queue,
	}
}

package study

import "main-service/internal/service"

type Implementation struct {
	studyService service.StudyService
}

func NewImplementation(studyService service.StudyService) *Implementation {
	return &Implementation{
		studyService: studyService,
	}
}

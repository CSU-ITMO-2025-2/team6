package study

import (
	"study-service/internal/service"
	desc "study-service/pkg/study_v1"
)

type Implementation struct {
	desc.UnimplementedStudyServiceServer
	studyService service.StudyService
}

func NewImplementation(studyService service.StudyService) *Implementation {
	return &Implementation{
		studyService: studyService,
	}
}

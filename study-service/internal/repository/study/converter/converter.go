package converter

import (
	"study-service/internal/model"
	modelRepo "study-service/internal/repository/study/model"
)

func ToStudyFromRepo(study *modelRepo.Study) *model.Study {
	m := &model.Study{
		ID:               study.ID,
		Name:             study.Name,
		Status:           model.StudyStatusType(study.Status),
		OwnerID:          study.OwnerID,
		ImageID:          study.ImageID,
		PredictedClassID: study.PredictedClassID,
		PredictedScore:   study.PredictedScore,
		ErrorDescription: study.ErrorDescription,
		CreatedAt:        study.CreatedAt,
		UpdatedAt:        study.UpdatedAt,
	}
	return m
}

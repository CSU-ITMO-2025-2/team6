package converter

import (
	"time"

	"main-service/internal/model"
	pb "study-service/pkg/study_v1"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func bytesToUUID(b []byte) uuid.UUID {
	if len(b) != 16 {
		return uuid.Nil
	}
	id, err := uuid.FromBytes(b)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func tsToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func PbStudyToModel(pbStudy *pb.StudyInfo) *model.Study {
	if pbStudy == nil {
		return nil
	}

	studyID := bytesToUUID(pbStudy.StudyId)
	ownerID := bytesToUUID(pbStudy.OwnerId)

	// predicted_score
	var scorePtr *float64
	if pbStudy.PredictedClassScore != nil {
		val := float64(*pbStudy.PredictedClassScore)
		scorePtr = &val
	}

	// predicted_class
	var classPtr *string
	if pbStudy.PredictedClass != nil {
		val := *pbStudy.PredictedClass
		classPtr = &val
	}

	// image URL
	var imageURL *string
	if pbStudy.ImageUrl != nil {
		val := *pbStudy.ImageUrl
		imageURL = &val
	}

	// name
	var namePtr *string
	if pbStudy.StudyName != nil {
		val := *pbStudy.StudyName
		namePtr = &val
	}

	// error description
	var errDesc *string
	if pbStudy.ErrorDescription != nil {
		val := *pbStudy.ErrorDescription
		errDesc = &val
	}

	createdAt := tsToTime(pbStudy.CreatedAt)
	updatedAt := pbStudy.UpdatedAt.AsTime()
	var updatedAtPtr *time.Time
	if !updatedAt.IsZero() {
		updatedAtPtr = &updatedAt
	}

	return &model.Study{
		ID:               studyID,
		Name:             namePtr,
		Status:           pbStudy.Status,
		OwnerID:          ownerID,
		ImageURL:         imageURL,
		PredictedClass:   classPtr,
		PredictedScore:   scorePtr,
		ErrorDescription: errDesc,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAtPtr,
	}
}

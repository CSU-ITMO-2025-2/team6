package study

import (
	"context"
	"fmt"

	desc "study-service/pkg/study_v1"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CreateStudy(ctx context.Context, req *desc.CreateStudyRequest) (*desc.CreateStudyResponse, error) {

	switch req.Image.MimeType {
	case "image/jpeg", "image/jpg", "image/png":
	default:
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid file type: %s", req.Image.MimeType))
	}

	if len(req.Image.ChunkData) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty file")
	}

	userID, _ := uuid.FromBytes(req.UserId)

	id, err := i.studyService.CreateStudy(ctx, userID, req.Image.ChunkData, req.Image.MimeType)
	if err != nil {
		return nil, err
	}

	return &desc.CreateStudyResponse{
		StudyId: id[:],
	}, nil
}

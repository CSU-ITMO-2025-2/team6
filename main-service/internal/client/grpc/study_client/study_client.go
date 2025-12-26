package study_client

import (
	"context"

	"main-service/internal/client/grpc/study_client/converter"
	"main-service/internal/model"
	pb "study-service/pkg/study_v1"

	"github.com/google/uuid"
)

type StudyClient struct {
	client pb.StudyServiceClient
}

func NewStudyClient(c pb.StudyServiceClient) *StudyClient {
	return &StudyClient{client: c}
}

func (c *StudyClient) CreateStudy(ctx context.Context, st *model.Study) (uuid.UUID, error) {
	resp, err := c.client.CreateStudy(ctx, &pb.CreateStudyRequest{
		UserId: st.OwnerID[:],
		Image: &pb.Image{
			ChunkData: st.Image.ChunkData,
			FileName:  st.Image.FileName,
			MimeType:  st.Image.MimeType,
		},
	})
	if err != nil {
		return uuid.Nil, err
	}

	studyID, err := uuid.FromBytes(resp.StudyId)
	if err != nil {
		return uuid.Nil, err
	}

	return studyID, nil
}

func (c *StudyClient) GetStudy(ctx context.Context, id uuid.UUID) (*model.Study, error) {
	resp, err := c.client.GetStudy(ctx, &pb.GetStudyRequest{
		StudyId: id[:],
	})
	if err != nil {
		return nil, err
	}

	return converter.PbStudyToModel(resp.GetStudy()), nil
}

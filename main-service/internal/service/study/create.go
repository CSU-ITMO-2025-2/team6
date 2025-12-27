package study

import (
	"context"

	"main-service/internal/model"

	"github.com/google/uuid"
)

func (s *serv) Create(ctx context.Context, st *model.Study) (uuid.UUID, error) {
	id, err := s.client.CreateStudy(ctx, st)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

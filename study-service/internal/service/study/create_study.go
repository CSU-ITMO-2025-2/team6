package study

import (
	"context"
	"fmt"

	"study-service/internal/model"

	"github.com/google/uuid"
)

func (s *serv) CreateStudy(ctx context.Context, userID uuid.UUID, image []byte, contentType string) (uuid.UUID, error) {

	// TODO подумать над транзакционным созданием
	imageID := uuid.New()
	if err := s.storage.Upload(ctx, userID, imageID, image, contentType); err != nil {
		return uuid.Nil, err
	}

	study := &model.Study{
		OwnerID: userID,
		ImageID: imageID,
		Status:  model.New,
	}

	id, err := s.studyRepository.Create(ctx, study)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("create study: %w", err)
	}

	msg := model.MLStream{
		StudyID: id,
		Image:   image,
	}

	if err := s.queue.PublishStreamML(msg); err != nil {
		return uuid.UUID{}, fmt.Errorf("publish stream ml: %w", err)
	}

	study.Status = model.Queued
	study.ID = id
	_, err = s.studyRepository.Update(ctx, study)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("update study status: %w", err)
	}

	// нужно создать исследование в БД
	// нужно закинуть фотку в миниО
	// нужно закинуть исследование в брокер
	//
	//

	return id, nil
}

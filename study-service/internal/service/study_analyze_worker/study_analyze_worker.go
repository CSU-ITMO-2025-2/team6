package study_analyze_worker

import (
	"context"
	"encoding/json"
	"time"

	"local-lib/queue"
	"study-service/internal/model"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const (
	imageAnalysisSubject = "core-task-ml.response"
	maxImageWorkers      = 30
)

type ImageAnalysisResponse struct {
	StudyID uuid.UUID `json:"study_id"`
	Error   *string   `json:"error,omitempty"`
	Type    string    `json:"type"`
	Score   int       `json:"score"`
}

func (s *serv) ImageAnalysisWorker(log *zap.SugaredLogger, n queue.Study) {
	pool := make(chan struct{}, maxImageWorkers)

	_, err := n.Subscribe(
		imageAnalysisSubject,
		func(msg *nats.Msg) {
			pool <- struct{}{}
			go func() {
				defer func() { <-pool }()
				s.handleImageAnalysisMessage(log, msg)
			}()
		},
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to image analysis responses: %s", err)
	}

	log.Info("Image Analysis Worker Started")
	select {}
}

func (s *serv) handleImageAnalysisMessage(log *zap.SugaredLogger, msg *nats.Msg) {
	var response ImageAnalysisResponse

	if err := json.Unmarshal(msg.Data, &response); err != nil {
		log.Errorf("Error unmarshalling image analysis response: %s", err)
		msg.NakWithDelay(5 * time.Second)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	study, err := s.studyRepository.Get(ctx, response.StudyID)
	if err != nil {
		log.Errorf("Error checking deduplication for study %s: %s", response.StudyID, err)
		msg.NakWithDelay(5 * time.Second)
		return
	}

	switch study.Status {
	case model.Completed, model.Failed:
		log.Infof("Study %s already processed, skipping", response.StudyID)
		msg.Ack()
		return
	}

	m := &model.Study{
		ID: response.StudyID,
	}

	if response.Error != nil {
		m.Status = model.Failed
		m.ErrorDescription = response.Error
	} else {
		m.Status = model.Completed
		m.PredictedScore = &response.Score
	}

	if _, err = s.studyRepository.Update(ctx, m); err != nil {
		log.Errorf("Error updating study: %s", err)
		msg.NakWithDelay(5 * time.Second)
		return
	}

	if err := msg.Ack(); err != nil {
		log.Errorf("Error ack message: %s", err)
	}

	log.Infof("Successfully processed image analysis for stduy: %s", response.StudyID)
}

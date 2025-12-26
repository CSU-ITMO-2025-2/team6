package study

import (
	"main-service/internal/client"
	"main-service/internal/service"
)

type serv struct {
	client client.StudyClient
}

func NewService(client client.StudyClient) service.StudyService {
	return &serv{
		client: client,
	}
}

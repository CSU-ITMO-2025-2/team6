package config

import (
	"log"
	"os"
)

const (
	storagePathEnvName = "STORAGE_PATH"
)

type ServiceConfig interface {
	StoragePath() string
}

type serviceConfig struct {
	storagePath string
}

func NewServiceConfig() (ServiceConfig, error) {
	storagePath := os.Getenv(storagePathEnvName)
	if storagePath == "" {
		storagePath = "./download"
		if err := os.MkdirAll(storagePath, 0755); err != nil {
			log.Fatalf("Failed to create storage directory: %v", err)
		}
	}

	return &serviceConfig{
		storagePath: storagePath,
	}, nil
}

func (c *serviceConfig) StoragePath() string {
	return c.storagePath
}

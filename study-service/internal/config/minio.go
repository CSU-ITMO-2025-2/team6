package config

import (
	"errors"
	"os"
)

const (
	minIOEnvHost      = "MINIO_HOST"
	minIOEnvUser      = "MINIO_ROOT_USER"
	minIOEnvPassword  = "MINIO_ROOT_PASSWORD"
	minIOEnvUseSSL    = "MINIO_USE_SSL"
	minIOEnvRemoteURL = "MINIO_REMOTE_URL"
)

type MinIOConfig interface {
	Host() string
	User() string
	Password() string
	UseSSL() bool
}

type minIOConfig struct {
	host      string
	user      string
	password  string
	useSSL    bool
	remoteURL string
}

func NewMinIOConfig() (MinIOConfig, error) {
	host := os.Getenv(minIOEnvHost)
	if len(host) == 0 {
		return nil, errors.New("minIO config error: host not found")
	}
	user := os.Getenv(minIOEnvUser)
	if len(user) == 0 {
		return nil, errors.New("minIO config error: user not found")
	}
	pass := os.Getenv(minIOEnvPassword)
	if len(pass) == 0 {
		return nil, errors.New("minIO config error: password not found")
	}
	useSSL := os.Getenv(minIOEnvUseSSL) == "true"
	if len(pass) == 0 {
		return nil, errors.New("minIO config error: useSSL not found")
	}
	remoteURL := os.Getenv(minIOEnvRemoteURL)
	if len(remoteURL) == 0 {
		return nil, errors.New("minIO config error: remoteURL not found")
	}

	return &minIOConfig{
		host:      host,
		user:      user,
		password:  pass,
		useSSL:    useSSL,
		remoteURL: remoteURL,
	}, nil
}

func (cfg *minIOConfig) Host() string {
	return cfg.host
}

func (cfg *minIOConfig) User() string {
	return cfg.user
}

func (cfg *minIOConfig) Password() string {
	return cfg.password
}

func (cfg *minIOConfig) UseSSL() bool {
	return cfg.useSSL
}

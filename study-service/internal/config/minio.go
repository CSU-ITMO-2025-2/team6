package config

import (
	"errors"
	"os"
)

const (
	s3EnvEndpoint  = "S3_ENDPOINT"
	s3EnvAccessKey = "S3_ACCESS_KEY"
	s3EnvSecretKey = "S3_SECRET_KEY"
	s3EnvBucket    = "S3_BUCKET"
	s3EnvUseSSL    = "S3_USE_SSL"
)

type S3Config interface {
	Endpoint() string
	AccessKey() string
	SecretKey() string
	Bucket() string
	UseSSL() bool
}

type s3Config struct {
	endpoint  string
	accessKey string
	secretKey string
	bucket    string
	useSSL    bool
}

func NewS3Config() (S3Config, error) {
	endpoint := os.Getenv(s3EnvEndpoint)
	if endpoint == "" {
		return nil, errors.New("s3 config error: endpoint not found")
	}

	accessKey := os.Getenv(s3EnvAccessKey)
	if accessKey == "" {
		return nil, errors.New("s3 config error: access key not found")
	}

	secretKey := os.Getenv(s3EnvSecretKey)
	if secretKey == "" {
		return nil, errors.New("s3 config error: secret key not found")
	}

	bucket := os.Getenv(s3EnvBucket)
	if bucket == "" {
		return nil, errors.New("s3 config error: bucket not found")
	}

	useSSL := os.Getenv(s3EnvUseSSL) == "true"

	return &s3Config{
		endpoint:  endpoint,
		accessKey: accessKey,
		secretKey: secretKey,
		bucket:    bucket,
		useSSL:    useSSL,
	}, nil
}

func (c *s3Config) Endpoint() string  { return c.endpoint }
func (c *s3Config) AccessKey() string { return c.accessKey }
func (c *s3Config) SecretKey() string { return c.secretKey }
func (c *s3Config) Bucket() string    { return c.bucket }
func (c *s3Config) UseSSL() bool      { return c.useSSL }

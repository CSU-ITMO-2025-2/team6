package minio

import (
	"bytes"
	"context"
	"io"

	"study-service/internal/config"
	"study-service/internal/repository"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type storage struct {
	client *minio.Client
	bucket string
}

func NewStorage(cfg config.S3Config) (repository.Storage, error) {
	client, err := minio.New(cfg.Endpoint(), &minio.Options{
		Creds: credentials.NewStaticV4(
			cfg.AccessKey(),
			cfg.SecretKey(),
			"",
		),
		Secure: cfg.UseSSL(),
	})
	if err != nil {
		return nil, err
	}

	return &storage{
		client: client,
		bucket: cfg.Bucket(),
	}, nil
}
func (s *storage) Upload(ctx context.Context, bucketName, objectName uuid.UUID, data []byte, contentType string) error {
	reader := bytes.NewReader(data)

	exists, err := s.client.BucketExists(ctx, bucketName.String())
	if err != nil {
		return err
	}
	if !exists {
		err := s.client.MakeBucket(ctx, bucketName.String(), minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	_, err = s.client.PutObject(ctx, bucketName.String(), objectName.String(), reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (s *storage) Download(ctx context.Context, bucketName, objectName string) (io.Reader, error) {
	object, err := s.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (s *storage) Delete(ctx context.Context, bucketName, objectName string) error {
	return s.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

//func (s *storage) GetURL(bucketName, objectName string) string {
//	// Для публичного доступа или через presigned URL
//	return s.endpoint + "/" + bucketName + "/" + objectName
//}

//func (s *storage) CreateOne(userID string, file helpers.FileDataType) (string, string, error) {
//	objectID := uuid.New().String()
//
//	ctx := context.Background()
//
//	exists, err := s.client.BucketExists(ctx, userID)
//	if err != nil {
//		return "", "", err
//	}
//	if !exists {
//		err := s.client.MakeBucket(ctx, userID, minio.MakeBucketOptions{})
//		if err != nil {
//			return "", "", err
//		}
//	}
//
//	reader := bytes.NewReader(file.Data)
//
//	res, err := s.client.PutObject(context.Background(), userID, objectID, reader, int64(len(file.Data)), minio.PutObjectOptions{})
//	if err != nil {
//		return "", "", fmt.Errorf("error creating object %s: %v", file.FileName, err)
//	}
//
//	url, err := s.client.PresignedGetObject(context.Background(), userID, objectID, time.Second*24*60*60, nil)
//	if err != nil {
//		return "", "", fmt.Errorf("ошибка при получении URL для объекта %s: %v", objectID, err)
//	}
//
//	return res.Key, url.String(), nil
//}
//
//func (s *storage) CreateMany(data map[string]helpers.FileDataType) ([]string, error) {
//	urls := make([]string, 0, len(data))
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	urlCh := make(chan string, len(data))
//
//	var wg sync.WaitGroup
//
//	for objectID, file := range data {
//		wg.Add(1)
//		go func(objectID string, file helpers.FileDataType) {
//			defer wg.Done()
//			_, err := m.mc.PutObject(ctx, m.BucketName, objectID, bytes.NewReader(file.Data), int64(len(file.Data)), minio.PutObjectOptions{})
//			if err != nil {
//				cancel()
//				return
//			}
//
//			url, err := m.mc.PresignedGetObject(ctx, m.BucketName, objectID, time.Second*24*60*60, nil)
//			if err != nil {
//				cancel()
//				return
//			}
//
//			urlCh <- url.String()
//		}(objectID, file)
//	}
//
//	go func() {
//		wg.Wait()
//		close(urlCh)
//	}()
//
//	for url := range urlCh {
//		urls = append(urls, url)
//	}
//
//	return urls, nil
//}
//
//func (s *storage) GetFileBytes(userID, objectID string) (*minio.Object, error) {
//	fileBytes, err := m.mc.GetObject(context.Background(), userID, objectID, minio.GetObjectOptions{})
//	if err != nil {
//		return nil, fmt.Errorf("error getting object %s: %v", objectID, err)
//	}
//
//	return fileBytes, nil
//}
//
//func (s *storage) GetFileUrl(userID, objectID string) (string, error) {
//	url, err := m.mc.PresignedGetObject(context.Background(), userID, objectID, time.Second*24*60*60, nil)
//	if err != nil {
//		return "", fmt.Errorf("ошибка при получении URL для объекта %s: %v", objectID, err)
//	}
//
//	// Костыль, чтобы работать с minio через nginx
//	fileURL := m.Endpoint + strings.TrimPrefix(url.String(), "http://82.97.241.8:9000")
//
//	return fileURL, nil
//}
//
//func (s *storage) GetMany(userID string, objectIDs []string) ([]string, error) {
//	fileCh := make(chan string, len(objectIDs))
//	errCh := make(chan helpers.OperationError, len(objectIDs))
//
//	var wg sync.WaitGroup
//	_, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	for _, objectID := range objectIDs {
//		wg.Add(1)
//		go func(objectID string) {
//			defer wg.Done()
//			url, err := s.GetFileUrl(userID, objectID)
//			if err != nil {
//				errCh <- helpers.OperationError{ObjectID: objectID, Error: fmt.Errorf("error getting object %s: %v", objectID, err)}
//				cancel()
//				return
//			}
//			fileCh <- url
//		}(objectID)
//	}
//
//	go func() {
//		wg.Wait()
//		close(fileCh)
//		close(errCh)
//	}()
//
//	var files []string
//	var errs []error
//	for url := range fileCh {
//		files = append(files, url)
//	}
//	for opErr := range errCh {
//		errs = append(errs, opErr.Error)
//	}
//
//	if len(errs) > 0 {
//		return nil, fmt.Errorf("error getting objects: %v", errs)
//	}
//
//	return files, nil
//}
//
//func (s *storage) DeleteOne(objectID string) error {
//	err := m.mc.RemoveObject(context.Background(), m.BucketName, objectID, minio.RemoveObjectOptions{})
//	if err != nil {
//		return err // Возвращаем ошибку, если не удалось удалить объект.
//	}
//	return nil
//}
//
//func (s *storage) DeleteMany(objectIDs []string) error {
//	errCh := make(chan helpers.OperationError, len(objectIDs))
//	var wg sync.WaitGroup
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	for _, objectID := range objectIDs {
//		wg.Add(1)
//		go func(id string) {
//			defer wg.Done()
//			err := m.mc.RemoveObject(ctx, m.BucketName, id, minio.RemoveObjectOptions{})
//			if err != nil {
//				errCh <- helpers.OperationError{ObjectID: id, Error: fmt.Errorf("error deleting objects %s: %v", id, err)}
//				cancel()
//			}
//		}(objectID)
//	}
//
//	go func() {
//		wg.Wait()
//		close(errCh)
//	}()
//
//	var errs []error
//	for opErr := range errCh {
//		errs = append(errs, opErr.Error)
//	}
//
//	if len(errs) > 0 {
//		return fmt.Errorf("error deleting objects : %v", errs) // Возврат ошибки, если возникли ошибки при удалении объектов
//	}
//
//	return nil
//}

package minio

import (
	"bytes"
	"context"
	"everlasting/src/infrastructure/pkg/logger"
	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	client *minio.Client
	logger *logger.AppLogger
}

func NewMinioStorage(client *minio.Client, logger *logger.AppLogger) *MinioStorage {
	return &MinioStorage{
		client: client,
		logger: logger,
	}
}

func (m *MinioStorage) Put(ctx context.Context, bucket, fileName string, file []byte, size int64, cacheAble bool, contentType string) error {
	cacheControl := "public, max-age=86400"
	if !cacheAble {
		cacheControl = "no-cache, max-age=0"
	}

	r := bytes.NewReader(file)

	_, err := m.client.PutObject(ctx, bucket, fileName, r, size, minio.PutObjectOptions{
		CacheControl: cacheControl,
		ContentType:  contentType,
	})

	if err != nil {
		m.logger.Error(ctx, "minio:put_error:", err.Error())
	}

	return err
}

func (m *MinioStorage) Get(ctx context.Context, bucket, fileName string) (file []byte, err error) {
	object, err := m.client.GetObject(ctx, bucket, fileName, minio.GetObjectOptions{})

	if err != nil {
		return nil, err
	}

	defer object.Close()

	_, err = object.Read(file)
	if err != nil {
		m.logger.Error(ctx, "minio:get_error", err.Error())
		return nil, err
	}

	return file, nil
}

func (m *MinioStorage) Delete(ctx context.Context, bucket, fileName string) error {
	err := m.client.RemoveObject(ctx, bucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		m.logger.Error(ctx, "minio:delete_error", err.Error())
		return err
	}

	return nil
}

func (m *MinioStorage) Exist(ctx context.Context, key string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

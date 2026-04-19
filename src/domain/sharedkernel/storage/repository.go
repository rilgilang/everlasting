package minio

import (
	"context"
)

//var (
//	ErrCacheKeyNotFound error = errors.New("cache key not found")
//)

type (
	StorageRepository interface {
		Put(ctx context.Context, bucket, fileName string, file []byte, size int64, cachable bool, contentType string) error
		Get(ctx context.Context, bucket, fileName string) (file []byte, err error)
		Delete(ctx context.Context, bucket, fileName string) error
		Exist(ctx context.Context, key string) (int64, error)
	}
)

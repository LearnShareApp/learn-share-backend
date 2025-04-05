package minio

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Host      string `env:"MINIO_HOST"       env-required:"true"`
	Port      int    `env:"MINIO_PORT"       env-required:"true"`
	AccessKey string `env:"MINIO_ACCESS_KEY" env-required:"true"`
	SecretKey string `env:"MINIO_SECRET_KEY" env-required:"true"`
	Bucket    string `env:"MINIO_BUCKET"     env-required:"true"`
	IsSSL     bool   `env:"IS_MINIO_SSL"     env-required:"true"`
}

func NewClient(config *Config) (*minio.Client, error) {
	minioClient, err := minio.New(
		fmt.Sprintf("%s:%d", config.Host, config.Port),
		&minio.Options{
			Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
			Secure: config.IsSSL,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to minIO: %w", err)
	}

	return minioClient, nil
}

func CreateBucket(ctx context.Context, client *minio.Client, bucketName string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return nil
}

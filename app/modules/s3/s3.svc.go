package s3

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"eduflow/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	UseSSL          bool
}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
}

type Service struct {
	tracer        trace.Tracer
	client        *minio.Client
	defaultBucket string
}

func newService(opt *Options) (*Service, error) {
	if opt == nil || opt.Val == nil {
		return nil, fmt.Errorf("%w", ErrS3InvalidConfig)
	}

	cfg := opt.Val
	if strings.TrimSpace(cfg.Endpoint) == "" || strings.TrimSpace(cfg.AccessKeyID) == "" || strings.TrimSpace(cfg.SecretAccessKey) == "" || strings.TrimSpace(cfg.BucketName) == "" {
		return &Service{tracer: opt.tracer, defaultBucket: strings.TrimSpace(cfg.BucketName)}, nil
	}

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrS3InvalidConfig, err)
	}

	return &Service{tracer: opt.tracer, client: client, defaultBucket: cfg.BucketName}, nil
}

func (s *Service) PresignUploadURL(ctx context.Context, bucket string, objectKey string, expires time.Duration) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("%w", ErrS3InvalidConfig)
	}
	if !isValidObjectKey(objectKey) {
		return "", fmt.Errorf("%w", ErrS3InvalidObjectKey)
	}
	bucketName, err := s.resolveBucket(bucket)
	if err != nil {
		return "", err
	}
	if expires <= 0 {
		expires = 15 * time.Minute
	}

	presignedURL, err := s.client.PresignedPutObject(ctx, bucketName, objectKey, expires)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrS3PresignFailed, err)
	}

	return presignedURL.String(), nil
}

func (s *Service) PresignDownloadURL(ctx context.Context, bucket string, objectKey string, expires time.Duration) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("%w", ErrS3InvalidConfig)
	}
	if !isValidObjectKey(objectKey) {
		return "", fmt.Errorf("%w", ErrS3InvalidObjectKey)
	}
	bucketName, err := s.resolveBucket(bucket)
	if err != nil {
		return "", err
	}
	if expires <= 0 {
		expires = 15 * time.Minute
	}

	presignedURL, err := s.client.PresignedGetObject(ctx, bucketName, objectKey, expires, url.Values{})
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrS3PresignFailed, err)
	}

	return presignedURL.String(), nil
}

func (s *Service) resolveBucket(bucket string) (string, error) {
	bucket = strings.TrimSpace(bucket)
	if bucket == "" {
		bucket = s.defaultBucket
	}
	if bucket == "" {
		return "", fmt.Errorf("%w", ErrS3InvalidConfig)
	}
	return bucket, nil
}

func isValidObjectKey(v string) bool {
	v = strings.TrimSpace(v)
	if v == "" || len(v) > 512 {
		return false
	}
	if strings.HasPrefix(v, "/") || strings.Contains(v, "..") {
		return false
	}
	return true
}

package s3

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
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
	Region          string
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
	endpoint := strings.TrimSpace(cfg.Endpoint)
	accessKey := strings.TrimSpace(cfg.AccessKeyID)
	secretKey := strings.TrimSpace(cfg.SecretAccessKey)
	bucketName := strings.TrimSpace(cfg.BucketName)
	region := strings.TrimSpace(cfg.Region)
	useSSL := cfg.UseSSL

	if endpoint == "" {
		endpoint = strings.TrimSpace(os.Getenv("OBJECT_ENDPOINT_URL"))
	}
	if accessKey == "" {
		accessKey = strings.TrimSpace(os.Getenv("OBJECT_ACCESS_KEY_ID"))
	}
	if secretKey == "" {
		secretKey = strings.TrimSpace(os.Getenv("OBJECT_SECRET_ACCESS_KEY"))
	}
	if bucketName == "" {
		bucketName = strings.TrimSpace(os.Getenv("OBJECT_PUBLIC_BUCKET"))
	}
	if region == "" {
		region = strings.TrimSpace(os.Getenv("OBJECT_REGION"))
	}

	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		if parsed, err := url.Parse(endpoint); err == nil {
			if parsed.Host != "" {
				endpoint = parsed.Host
			}
			if parsed.Scheme == "https" {
				useSSL = true
			}
		}
	}

	if endpoint == "" || accessKey == "" || secretKey == "" || bucketName == "" {
		return &Service{tracer: opt.tracer, defaultBucket: bucketName}, nil
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure:       useSSL,
		Region:       region,
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrS3InvalidConfig, err)
	}

	return &Service{tracer: opt.tracer, client: client, defaultBucket: bucketName}, nil
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

func (s *Service) UploadObject(ctx context.Context, bucket string, objectKey string, reader io.Reader, size int64, contentType string) error {
	if s.client == nil {
		return fmt.Errorf("%w", ErrS3InvalidConfig)
	}
	if !isValidObjectKey(objectKey) {
		return fmt.Errorf("%w", ErrS3InvalidObjectKey)
	}
	bucketName, err := s.resolveBucket(bucket)
	if err != nil {
		return err
	}
	if reader == nil {
		return fmt.Errorf("%w", ErrS3UploadFailed)
	}
	if size < 0 {
		return fmt.Errorf("%w", ErrS3UploadFailed)
	}
	if strings.TrimSpace(contentType) == "" {
		contentType = "application/octet-stream"
	}

	_, err = s.client.PutObject(ctx, bucketName, objectKey, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrS3UploadFailed, err)
	}

	return nil
}

func isValidObjectKey(v string) bool {
	v = strings.TrimSpace(v)
	if v == "" || len(v) > 512 {
		return false
	}
	if strings.HasPrefix(v, "/") {
		return false
	}

	normalized := strings.ReplaceAll(v, "\\", "/")
	if strings.HasPrefix(normalized, "../") || strings.HasSuffix(normalized, "/..") || strings.Contains(normalized, "/../") {
		return false
	}
	return true
}

package pictures

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/modules/s3"
	"eduflow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	DefaultUploadExpireSeconds   int64
	DefaultDownloadExpireSeconds int64
	DefaultBucket                string
}

type Options struct {
	*config.Config[Config]
	tracer    trace.Tracer
	s3        *s3.Service
	db        entitiesinf.DocumentEntity
	storageDB entitiesinf.StorageEntity
}

type Service struct {
	tracer                trace.Tracer
	s3                    *s3.Service
	db                    entitiesinf.DocumentEntity
	storageDB             entitiesinf.StorageEntity
	defaultUploadExpire   time.Duration
	defaultDownloadExpire time.Duration
	defaultBucket         string
}

func newService(opt *Options) *Service {
	uploadExpire := time.Duration(opt.Val.DefaultUploadExpireSeconds) * time.Second
	if uploadExpire <= 0 {
		uploadExpire = 15 * time.Minute
	}

	downloadExpire := time.Duration(opt.Val.DefaultDownloadExpireSeconds) * time.Second
	if downloadExpire <= 0 {
		downloadExpire = 15 * time.Minute
	}

	defaultBucket := strings.TrimSpace(opt.Val.DefaultBucket)
	if defaultBucket == "" {
		defaultBucket = strings.TrimSpace(os.Getenv("PICTURES__DEFAULT_BUCKET"))
	}
	if defaultBucket == "" {
		defaultBucket = strings.TrimSpace(os.Getenv("OBJECT_PUBLIC_BUCKET"))
	}

	return &Service{
		tracer:                opt.tracer,
		s3:                    opt.s3,
		db:                    opt.db,
		storageDB:             opt.storageDB,
		defaultUploadExpire:   uploadExpire,
		defaultDownloadExpire: downloadExpire,
		defaultBucket:         defaultBucket,
	}
}

func parsePictureStatus(v string) (ent.DocumentStatus, bool) {
	s := ent.DocumentStatus(strings.TrimSpace(v))
	switch s {
	case ent.DocumentStatusPendingUpload, ent.DocumentStatusActive, ent.DocumentStatusArchived:
		return s, true
	default:
		return "", false
	}
}

func (s *Service) resolvePictureStorage(ctx context.Context, schoolID uuid.UUID) (*ent.Storage, error) {
	if s.defaultBucket == "" {
		return nil, fmt.Errorf("%w", ErrPictureConditionFail)
	}

	item, err := s.storageDB.EnsureStorageByBucket(ctx, schoolID, s.defaultBucket, nil)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

func (s *Service) PresignUploadURL(ctx context.Context, objectKey string, expiresSeconds *int64) (string, error) {
	expires := s.defaultUploadExpire
	if expiresSeconds != nil && *expiresSeconds > 0 {
		expires = time.Duration(*expiresSeconds) * time.Second
	}
	if expires > 7*24*time.Hour {
		return "", fmt.Errorf("%w", ErrPictureConditionFail)
	}

	url, err := s.s3.PresignUploadURL(ctx, s.defaultBucket, objectKey, expires)
	if err != nil {
		return "", normalizeServiceError(err)
	}

	return url, nil
}

func (s *Service) PresignDownloadURL(ctx context.Context, objectKey string, expiresSeconds *int64) (string, error) {
	expires := s.defaultDownloadExpire
	if expiresSeconds != nil && *expiresSeconds > 0 {
		expires = time.Duration(*expiresSeconds) * time.Second
	}
	if expires > 7*24*time.Hour {
		return "", fmt.Errorf("%w", ErrPictureConditionFail)
	}

	url, err := s.s3.PresignDownloadURL(ctx, s.defaultBucket, objectKey, expires)
	if err != nil {
		return "", normalizeServiceError(err)
	}

	return url, nil
}

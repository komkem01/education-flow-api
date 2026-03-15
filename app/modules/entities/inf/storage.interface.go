package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type StorageEntity interface {
	GetStorageByID(ctx context.Context, id uuid.UUID, schoolID uuid.UUID) (*ent.Storage, error)
	GetStorageByBucket(ctx context.Context, schoolID uuid.UUID, bucketName string) (*ent.Storage, error)
	ListStorages(ctx context.Context, req *base.RequestPaginate, schoolID uuid.UUID, provider *ent.StorageProvider, isDefault *bool) ([]*ent.Storage, *base.ResponsePaginate, error)
	CreateStorage(ctx context.Context, data *ent.Storage) (*ent.Storage, error)
	UpdateStorageByID(ctx context.Context, id uuid.UUID, schoolID uuid.UUID, data *ent.StorageUpdate) (*ent.Storage, error)
	EnsureStorageByBucket(ctx context.Context, schoolID uuid.UUID, bucketName string, endpoint *string) (*ent.Storage, error)
}

package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.StorageEntity = (*Service)(nil)

func (s *Service) GetStorageByID(ctx context.Context, id uuid.UUID, schoolID uuid.UUID) (*ent.Storage, error) {
	item := new(ent.Storage)
	if err := s.db.NewSelect().
		Model(item).
		Where("id = ?", id).
		Where("school_id = ?", schoolID).
		Scan(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) GetStorageByBucket(ctx context.Context, schoolID uuid.UUID, bucketName string) (*ent.Storage, error) {
	item := new(ent.Storage)
	if err := s.db.NewSelect().
		Model(item).
		Where("school_id = ?", schoolID).
		Where("bucket_name = ?", bucketName).
		Scan(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) ListStorages(ctx context.Context, req *base.RequestPaginate, schoolID uuid.UUID, provider *ent.StorageProvider, isDefault *bool) ([]*ent.Storage, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.Storage, 0)
	query := s.db.NewSelect().Model(&items).Where("school_id = ?", schoolID)

	if provider != nil {
		query.Where("provider = ?", *provider)
	}
	if isDefault != nil {
		query.Where("is_default = ?", *isDefault)
	}

	if err := req.SetSearchBy(query, []string{"name", "bucket_name", "provider", "endpoint"}); err != nil {
		return nil, nil, err
	}
	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "updated_at", "name", "bucket_name", "provider", "is_default"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) CreateStorage(ctx context.Context, data *ent.Storage) (*ent.Storage, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) UpdateStorageByID(ctx context.Context, id uuid.UUID, schoolID uuid.UUID, data *ent.StorageUpdate) (*ent.Storage, error) {
	query := s.db.NewUpdate().Model(&ent.Storage{}).Where("id = ?", id).Where("school_id = ?", schoolID).Set("updated_at = now()")

	if data.Provider != nil {
		query.Set("provider = ?", *data.Provider)
	}
	if data.Name != nil {
		query.Set("name = ?", *data.Name)
	}
	if data.Endpoint != nil {
		query.Set("endpoint = ?", *data.Endpoint)
	}
	if data.BucketName != nil {
		query.Set("bucket_name = ?", *data.BucketName)
	}
	if data.IsDefault != nil {
		query.Set("is_default = ?", *data.IsDefault)
	}
	if data.Config != nil {
		query.Set("config = ?", *data.Config)
	}

	res, err := query.Exec(ctx)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, sql.ErrNoRows
	}

	return s.GetStorageByID(ctx, id, schoolID)
}

func (s *Service) EnsureStorageByBucket(ctx context.Context, schoolID uuid.UUID, bucketName string, endpoint *string) (*ent.Storage, error) {
	bucketName = strings.TrimSpace(bucketName)
	if bucketName == "" {
		return nil, fmt.Errorf("storage-bucket-required")
	}

	item, err := s.GetStorageByBucket(ctx, schoolID, bucketName)
	if err == nil {
		return item, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	create := &ent.Storage{
		SchoolID:   schoolID,
		Provider:   ent.StorageProviderS3,
		Name:       bucketName,
		Endpoint:   endpoint,
		BucketName: bucketName,
		IsDefault:  false,
	}

	created, createErr := s.CreateStorage(ctx, create)
	if createErr == nil {
		return created, nil
	}

	// Handle race on unique constraint by fetching existing record.
	fallback, getErr := s.GetStorageByBucket(ctx, schoolID, bucketName)
	if getErr == nil {
		return fallback, nil
	}

	return nil, createErr
}

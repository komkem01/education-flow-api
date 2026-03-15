package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.DocumentEntity = (*Service)(nil)

func (s *Service) CreateDocument(ctx context.Context, data *ent.Document) (*ent.Document, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetDocumentByID(ctx context.Context, id uuid.UUID, schoolID uuid.UUID) (*ent.Document, error) {
	item := new(ent.Document)
	if err := s.db.NewSelect().
		Model(item).
		ColumnExpr("st.name AS storage_name").
		Join("JOIN storages AS st ON st.id = d.storage_id").
		Where("d.id = ?", id).
		Where("d.school_id = ?", schoolID).
		Scan(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) ListDocuments(ctx context.Context, req *base.RequestPaginate, schoolID uuid.UUID, ownerMemberID *uuid.UUID, status *ent.DocumentStatus, storageID *uuid.UUID) ([]*ent.Document, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.Document, 0)
	query := s.db.NewSelect().
		Model(&items).
		ColumnExpr("st.name AS storage_name").
		Join("JOIN storages AS st ON st.id = d.storage_id").
		Where("d.school_id = ?", schoolID)

	if ownerMemberID != nil {
		query.Where("d.owner_member_id = ?", *ownerMemberID)
	}
	if status != nil {
		query.Where("d.status = ?", *status)
	}
	if storageID != nil {
		query.Where("d.storage_id = ?", *storageID)
	}

	if err := req.SetSearchBy(query, []string{"file_name", "object_key", "content_type", "storage_id"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "updated_at", "status", "size_bytes", "file_name", "storage_id"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateDocumentByID(ctx context.Context, id uuid.UUID, schoolID uuid.UUID, data *ent.DocumentUpdate) (*ent.Document, error) {
	query := s.db.NewUpdate().Model(&ent.Document{}).Where("id = ?", id).Where("school_id = ?", schoolID).Set("updated_at = now()")

	if data.OwnerMemberID != nil {
		query.Set("owner_member_id = ?", *data.OwnerMemberID)
	}
	if data.FileName != nil {
		query.Set("file_name = ?", *data.FileName)
	}
	if data.ContentType != nil {
		query.Set("content_type = ?", *data.ContentType)
	}
	if data.SizeBytes != nil {
		query.Set("size_bytes = ?", *data.SizeBytes)
	}
	if data.Status != nil {
		query.Set("status = ?", *data.Status)
	}
	if data.Metadata != nil {
		query.Set("metadata = ?", *data.Metadata)
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

	return s.GetDocumentByID(ctx, id, schoolID)
}

func (s *Service) SoftDeleteDocumentByID(ctx context.Context, id uuid.UUID, schoolID uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.Document{}).Where("id = ?", id).Where("school_id = ?", schoolID).Exec(ctx)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

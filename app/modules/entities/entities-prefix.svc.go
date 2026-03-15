package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.PrefixEntity = (*Service)(nil)

func (s *Service) CreatePrefix(ctx context.Context, genderID uuid.UUID, name string, isActive bool) (*ent.Prefix, error) {
	prefix := &ent.Prefix{
		GenderID: genderID,
		Name:     name,
		IsActive: isActive,
	}

	if _, err := s.db.NewInsert().Model(prefix).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return prefix, nil
}

func (s *Service) GetPrefixByID(ctx context.Context, id uuid.UUID) (*ent.Prefix, error) {
	prefix := new(ent.Prefix)
	if err := s.db.NewSelect().Model(prefix).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return prefix, nil
}

func (s *Service) ListPrefixes(ctx context.Context, req *base.RequestPaginate, isActive *bool, genderID *uuid.UUID) ([]*ent.Prefix, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	prefixes := make([]*ent.Prefix, 0)
	query := s.db.NewSelect().Model(&prefixes)

	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}

	if genderID != nil {
		query.Where("gender_id = ?", *genderID)
	}

	if err := req.SetSearchBy(query, []string{"name"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}

	if err := req.SetSortOrder(query, []string{"created_at", "name", "is_active", "gender_id"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)

	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return prefixes, &base.ResponsePaginate{
		Page:  req.GetPage(),
		Size:  req.GetSize(),
		Total: int64(total),
	}, nil
}

func (s *Service) UpdatePrefixByID(ctx context.Context, id uuid.UUID, genderID *uuid.UUID, name *string, isActive *bool) (*ent.Prefix, error) {
	query := s.db.NewUpdate().
		Model(&ent.Prefix{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if genderID != nil {
		query.Set("gender_id = ?", *genderID)
	}

	if name != nil {
		query.Set("name = ?", *name)
	}

	if isActive != nil {
		query.Set("is_active = ?", *isActive)
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

	return s.GetPrefixByID(ctx, id)
}

func (s *Service) SoftDeletePrefixByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.Prefix{}).
			Set("is_active = ?", false).
			Set("updated_at = now()").
			Where("id = ?", id).
			Exec(ctx)
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

		_, err = tx.NewDelete().Model(&ent.Prefix{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}

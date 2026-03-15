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

var _ entitiesinf.GenderEntity = (*Service)(nil)

func (s *Service) CreateGender(ctx context.Context, name string, isActive bool) (*ent.Gender, error) {
	gender := &ent.Gender{
		Name:     name,
		IsActive: isActive,
	}

	if _, err := s.db.NewInsert().Model(gender).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return gender, nil
}

func (s *Service) GetGenderByID(ctx context.Context, id uuid.UUID) (*ent.Gender, error) {
	gender := new(ent.Gender)
	if err := s.db.NewSelect().Model(gender).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return gender, nil
}

func (s *Service) ListGenders(ctx context.Context, req *base.RequestPaginate, isActive *bool) ([]*ent.Gender, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	genders := make([]*ent.Gender, 0)
	query := s.db.NewSelect().Model(&genders)

	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}

	if err := req.SetSearchBy(query, []string{"name"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}

	if err := req.SetSortOrder(query, []string{"created_at", "name", "is_active"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)

	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return genders, &base.ResponsePaginate{
		Page:  req.GetPage(),
		Size:  req.GetSize(),
		Total: int64(total),
	}, nil
}

func (s *Service) UpdateGenderByID(ctx context.Context, id uuid.UUID, name *string, isActive *bool) (*ent.Gender, error) {
	query := s.db.NewUpdate().
		Model(&ent.Gender{}).
		Where("id = ?", id).
		Set("updated_at = now()")

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

	return s.GetGenderByID(ctx, id)
}

func (s *Service) SoftDeleteGenderByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.Gender{}).
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

		_, err = tx.NewDelete().Model(&ent.Gender{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}

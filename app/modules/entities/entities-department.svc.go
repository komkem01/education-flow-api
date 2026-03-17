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

var _ entitiesinf.DepartmentEntity = (*Service)(nil)

func (s *Service) CreateDepartment(ctx context.Context, code string, name string, isActive bool) (*ent.Department, error) {
	department := &ent.Department{
		Code:     code,
		Name:     name,
		IsActive: isActive,
	}

	if _, err := s.db.NewInsert().Model(department).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return department, nil
}

func (s *Service) GetDepartmentByID(ctx context.Context, id uuid.UUID) (*ent.Department, error) {
	department := new(ent.Department)
	if err := s.db.NewSelect().Model(department).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return department, nil
}

func (s *Service) ListDepartments(ctx context.Context, req *base.RequestPaginate, isActive *bool) ([]*ent.Department, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	departments := make([]*ent.Department, 0)
	query := s.db.NewSelect().Model(&departments)

	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}

	if err := req.SetSearchBy(query, []string{"code", "name"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}

	if err := req.SetSortOrder(query, []string{"created_at", "code", "name", "is_active"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)

	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return departments, &base.ResponsePaginate{
		Page:  req.GetPage(),
		Size:  req.GetSize(),
		Total: int64(total),
	}, nil
}

func (s *Service) UpdateDepartmentByID(ctx context.Context, id uuid.UUID, code *string, name *string, isActive *bool) (*ent.Department, error) {
	query := s.db.NewUpdate().
		Model(&ent.Department{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if code != nil {
		query.Set("code = ?", *code)
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

	return s.GetDepartmentByID(ctx, id)
}

func (s *Service) SoftDeleteDepartmentByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.Department{}).
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

		_, err = tx.NewDelete().Model(&ent.Department{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}

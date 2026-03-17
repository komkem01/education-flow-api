package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.SchoolDepartmentEntity = (*Service)(nil)

func (s *Service) CreateSchoolDepartment(ctx context.Context, data *ent.SchoolDepartment) (*ent.SchoolDepartment, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetSchoolDepartmentByID(ctx context.Context, id uuid.UUID) (*ent.SchoolDepartment, error) {
	row := new(ent.SchoolDepartment)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) GetSchoolDepartmentBySchoolAndDepartment(ctx context.Context, schoolID uuid.UUID, departmentID uuid.UUID) (*ent.SchoolDepartment, error) {
	row := new(ent.SchoolDepartment)
	if err := s.db.NewSelect().
		Model(row).
		Where("school_id = ?", schoolID).
		Where("department_id = ?", departmentID).
		Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListSchoolDepartments(ctx context.Context, req *base.RequestPaginate, schoolID *uuid.UUID, departmentID *uuid.UUID, isActive *bool) ([]*ent.SchoolDepartment, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.SchoolDepartment, 0)
	query := s.db.NewSelect().Model(&items)

	if schoolID != nil {
		query.Where("school_id = ?", *schoolID)
	}
	if departmentID != nil {
		query.Where("department_id = ?", *departmentID)
	}
	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}

	if err := req.SetSearchBy(query, []string{"code", "custom_name"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "updated_at", "code", "is_active"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	count, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(count)}, nil
}

func (s *Service) UpdateSchoolDepartmentByID(ctx context.Context, id uuid.UUID, data *ent.SchoolDepartmentUpdate) (*ent.SchoolDepartment, error) {
	query := s.db.NewUpdate().
		Model(&ent.SchoolDepartment{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.SchoolID != nil {
		query.Set("school_id = ?", *data.SchoolID)
	}
	if data.DepartmentID != nil {
		query.Set("department_id = ?", *data.DepartmentID)
	}
	if data.Code != nil {
		query.Set("code = ?", *data.Code)
	}
	if data.CustomName != nil {
		query.Set("custom_name = ?", *data.CustomName)
	}
	if data.IsActive != nil {
		query.Set("is_active = ?", *data.IsActive)
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

	return s.GetSchoolDepartmentByID(ctx, id)
}

func (s *Service) SoftDeleteSchoolDepartmentByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.SchoolDepartment{}).Where("id = ?", id).Exec(ctx)
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

package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.TeacherLicenseEntity = (*Service)(nil)

func (s *Service) CreateTeacherLicense(ctx context.Context, data *ent.TeacherLicense) (*ent.TeacherLicense, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetTeacherLicenseByID(ctx context.Context, id uuid.UUID) (*ent.TeacherLicense, error) {
	row := new(ent.TeacherLicense)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListTeacherLicenses(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, status *ent.TeacherLicenseStatus) ([]*ent.TeacherLicense, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.TeacherLicense, 0)
	query := s.db.NewSelect().Model(&items)

	if teacherID != nil {
		query.Where("teacher_id = ?", *teacherID)
	}
	if status != nil {
		query.Where("license_status = ?", *status)
	}

	if err := req.SetSearchBy(query, []string{"license_no", "issued_by", "note"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "license_no", "issued_at", "expires_at", "license_status"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateTeacherLicenseByID(ctx context.Context, id uuid.UUID, data *ent.TeacherLicenseUpdate) (*ent.TeacherLicense, error) {
	query := s.db.NewUpdate().
		Model(&ent.TeacherLicense{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.TeacherID != nil {
		query.Set("teacher_id = ?", *data.TeacherID)
	}
	if data.LicenseNo != nil {
		query.Set("license_no = ?", *data.LicenseNo)
	}
	if data.IssuedAt != nil {
		query.Set("issued_at = ?", *data.IssuedAt)
	}
	if data.ExpiresAt != nil {
		query.Set("expires_at = ?", *data.ExpiresAt)
	}
	if data.LicenseStatus != nil {
		query.Set("license_status = ?", *data.LicenseStatus)
	}
	if data.IssuedBy != nil {
		query.Set("issued_by = ?", *data.IssuedBy)
	}
	if data.Note != nil {
		query.Set("note = ?", *data.Note)
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

	return s.GetTeacherLicenseByID(ctx, id)
}

func (s *Service) SoftDeleteTeacherLicenseByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.TeacherLicense{}).Where("id = ?", id).Exec(ctx)
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

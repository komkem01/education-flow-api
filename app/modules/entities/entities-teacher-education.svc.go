package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.TeacherEducationEntity = (*Service)(nil)

func (s *Service) CreateTeacherEducation(ctx context.Context, data *ent.TeacherEducation) (*ent.TeacherEducation, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetTeacherEducationByID(ctx context.Context, id uuid.UUID) (*ent.TeacherEducation, error) {
	row := new(ent.TeacherEducation)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListTeacherEducations(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, degree *ent.TeacherDegree) ([]*ent.TeacherEducation, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.TeacherEducation, 0)
	query := s.db.NewSelect().Model(&items)

	if teacherID != nil {
		query.Where("teacher_id = ?", *teacherID)
	}
	if degree != nil {
		query.Where("degree = ?", *degree)
	}

	if err := req.SetSearchBy(query, []string{"major", "university", "graduation_year"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "degree", "graduation_year", "major"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateTeacherEducationByID(ctx context.Context, id uuid.UUID, data *ent.TeacherEducationUpdate) (*ent.TeacherEducation, error) {
	query := s.db.NewUpdate().
		Model(&ent.TeacherEducation{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.TeacherID != nil {
		query.Set("teacher_id = ?", *data.TeacherID)
	}
	if data.Degree != nil {
		query.Set("degree = ?", *data.Degree)
	}
	if data.Major != nil {
		query.Set("major = ?", *data.Major)
	}
	if data.University != nil {
		query.Set("university = ?", *data.University)
	}
	if data.GraduationYear != nil {
		query.Set("graduation_year = ?", *data.GraduationYear)
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

	return s.GetTeacherEducationByID(ctx, id)
}

func (s *Service) SoftDeleteTeacherEducationByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.TeacherEducation{}).Where("id = ?", id).Exec(ctx)
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

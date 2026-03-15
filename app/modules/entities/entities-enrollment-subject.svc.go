package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.EnrollmentSubjectEntity = (*Service)(nil)

func (s *Service) CreateEnrollmentSubject(ctx context.Context, data *ent.EnrollmentSubject) (*ent.EnrollmentSubject, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetEnrollmentSubjectByID(ctx context.Context, id uuid.UUID) (*ent.EnrollmentSubject, error) {
	row := new(ent.EnrollmentSubject)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListEnrollmentSubjects(ctx context.Context, req *base.RequestPaginate, enrollmentID *uuid.UUID, subjectID *uuid.UUID, teacherID *uuid.UUID, status *ent.StudentEnrollmentStatus, isPrimary *bool) ([]*ent.EnrollmentSubject, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.EnrollmentSubject, 0)
	query := s.db.NewSelect().Model(&items)

	if enrollmentID != nil {
		query.Where("enrollment_id = ?", *enrollmentID)
	}
	if subjectID != nil {
		query.Where("subject_id = ?", *subjectID)
	}
	if teacherID != nil {
		query.Where("teacher_id = ?", *teacherID)
	}
	if status != nil {
		query.Where("status = ?", *status)
	}
	if isPrimary != nil {
		query.Where("is_primary = ?", *isPrimary)
	}

	if err := req.SetSearchBy(query, []string{"status"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "status", "is_primary"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateEnrollmentSubjectByID(ctx context.Context, id uuid.UUID, data *ent.EnrollmentSubjectUpdate) (*ent.EnrollmentSubject, error) {
	query := s.db.NewUpdate().
		Model(&ent.EnrollmentSubject{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.EnrollmentID != nil {
		query.Set("enrollment_id = ?", *data.EnrollmentID)
	}
	if data.SubjectID != nil {
		query.Set("subject_id = ?", *data.SubjectID)
	}
	if data.TeacherID != nil {
		query.Set("teacher_id = ?", *data.TeacherID)
	}
	if data.IsPrimary != nil {
		query.Set("is_primary = ?", *data.IsPrimary)
	}
	if data.Status != nil {
		query.Set("status = ?", *data.Status)
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

	return s.GetEnrollmentSubjectByID(ctx, id)
}

func (s *Service) SoftDeleteEnrollmentSubjectByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.EnrollmentSubject{}).Where("id = ?", id).Exec(ctx)
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

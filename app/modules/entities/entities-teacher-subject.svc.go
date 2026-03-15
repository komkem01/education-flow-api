package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.TeacherSubjectEntity = (*Service)(nil)

func (s *Service) CreateTeacherSubject(ctx context.Context, data *ent.TeacherSubject) (*ent.TeacherSubject, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetTeacherSubjectByID(ctx context.Context, id uuid.UUID) (*ent.TeacherSubject, error) {
	row := new(ent.TeacherSubject)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListTeacherSubjects(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, subjectID *uuid.UUID, role *ent.TeacherSubjectRole, isActive *bool) ([]*ent.TeacherSubject, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.TeacherSubject, 0)
	query := s.db.NewSelect().Model(&items)

	if teacherID != nil {
		query.Where("teacher_id = ?", *teacherID)
	}
	if subjectID != nil {
		query.Where("subject_id = ?", *subjectID)
	}
	if role != nil {
		query.Where("role = ?", *role)
	}
	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}

	if err := req.SetSearchBy(query, []string{"role"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "role", "is_active"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateTeacherSubjectByID(ctx context.Context, id uuid.UUID, data *ent.TeacherSubjectUpdate) (*ent.TeacherSubject, error) {
	query := s.db.NewUpdate().
		Model(&ent.TeacherSubject{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.TeacherID != nil {
		query.Set("teacher_id = ?", *data.TeacherID)
	}
	if data.SubjectID != nil {
		query.Set("subject_id = ?", *data.SubjectID)
	}
	if data.Role != nil {
		query.Set("role = ?", *data.Role)
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

	return s.GetTeacherSubjectByID(ctx, id)
}

func (s *Service) SoftDeleteTeacherSubjectByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.TeacherSubject{}).Where("id = ?", id).Exec(ctx)
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

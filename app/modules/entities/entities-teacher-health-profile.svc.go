package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.TeacherHealthProfileEntity = (*Service)(nil)

func (s *Service) CreateTeacherHealthProfile(ctx context.Context, data *ent.TeacherHealthProfile) (*ent.TeacherHealthProfile, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetTeacherHealthProfileByID(ctx context.Context, id uuid.UUID) (*ent.TeacherHealthProfile, error) {
	row := new(ent.TeacherHealthProfile)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListTeacherHealthProfiles(ctx context.Context, req *base.RequestPaginate, memberTeacherID *uuid.UUID, bloodType *string) ([]*ent.TeacherHealthProfile, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.TeacherHealthProfile, 0)
	query := s.db.NewSelect().Model(&items)

	if memberTeacherID != nil {
		query.Where("member_teacher_id = ?", *memberTeacherID)
	}
	if bloodType != nil {
		query.Where("blood_type = ?", *bloodType)
	}

	if err := req.SetSearchBy(query, []string{"blood_type", "allergy_info", "chronic_disease", "medication_note", "fitness_for_work_note"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "blood_type"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateTeacherHealthProfileByID(ctx context.Context, id uuid.UUID, data *ent.TeacherHealthProfileUpdate) (*ent.TeacherHealthProfile, error) {
	query := s.db.NewUpdate().
		Model(&ent.TeacherHealthProfile{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.MemberTeacherID != nil {
		query.Set("member_teacher_id = ?", *data.MemberTeacherID)
	}
	if data.BloodType != nil {
		query.Set("blood_type = ?", *data.BloodType)
	}
	if data.AllergyInfo != nil {
		query.Set("allergy_info = ?", *data.AllergyInfo)
	}
	if data.ChronicDisease != nil {
		query.Set("chronic_disease = ?", *data.ChronicDisease)
	}
	if data.MedicationNote != nil {
		query.Set("medication_note = ?", *data.MedicationNote)
	}
	if data.FitnessForWorkNote != nil {
		query.Set("fitness_for_work_note = ?", *data.FitnessForWorkNote)
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

	return s.GetTeacherHealthProfileByID(ctx, id)
}

func (s *Service) SoftDeleteTeacherHealthProfileByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.TeacherHealthProfile{}).Where("id = ?", id).Exec(ctx)
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

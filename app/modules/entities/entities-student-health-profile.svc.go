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

var _ entitiesinf.StudentHealthProfileEntity = (*Service)(nil)

func (s *Service) CreateStudentHealthProfile(ctx context.Context, data *ent.StudentHealthProfile) (*ent.StudentHealthProfile, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetStudentHealthProfileByID(ctx context.Context, id uuid.UUID) (*ent.StudentHealthProfile, error) {
	row := new(ent.StudentHealthProfile)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListStudentHealthProfiles(ctx context.Context, req *base.RequestPaginate, studentID *uuid.UUID, bloodType *string) ([]*ent.StudentHealthProfile, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.StudentHealthProfile, 0)
	query := s.db.NewSelect().Model(&items)

	if studentID != nil {
		query.Where("student_id = ?", *studentID)
	}
	if bloodType != nil {
		query.Where("blood_type = ?", *bloodType)
	}

	if err := req.SetSearchBy(query, []string{"blood_type", "allergy_info", "chronic_disease", "medical_note"}); err != nil {
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

func (s *Service) UpdateStudentHealthProfileByID(ctx context.Context, id uuid.UUID, data *ent.StudentHealthProfileUpdate) (*ent.StudentHealthProfile, error) {
	query := s.db.NewUpdate().
		Model(&ent.StudentHealthProfile{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.StudentID != nil {
		query.Set("student_id = ?", *data.StudentID)
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
	if data.MedicalNote != nil {
		query.Set("medical_note = ?", *data.MedicalNote)
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

	return s.GetStudentHealthProfileByID(ctx, id)
}

func (s *Service) SoftDeleteStudentHealthProfileByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.StudentHealthProfile{}).
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

		_, err = tx.NewDelete().Model(&ent.StudentHealthProfile{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}

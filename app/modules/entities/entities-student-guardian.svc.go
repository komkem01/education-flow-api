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

var _ entitiesinf.StudentGuardianEntity = (*Service)(nil)

func (s *Service) CreateStudentGuardian(ctx context.Context, data *ent.StudentGuardian) (*ent.StudentGuardian, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetStudentGuardianByID(ctx context.Context, id uuid.UUID) (*ent.StudentGuardian, error) {
	row := new(ent.StudentGuardian)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListStudentGuardians(ctx context.Context, req *base.RequestPaginate, studentID *uuid.UUID, guardianID *uuid.UUID, isMainGuardian *bool) ([]*ent.StudentGuardian, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.StudentGuardian, 0)
	query := s.db.NewSelect().Model(&items)

	if studentID != nil {
		query.Where("student_id = ?", *studentID)
	}
	if guardianID != nil {
		query.Where("guardian_id = ?", *guardianID)
	}
	if isMainGuardian != nil {
		query.Where("is_main_guardian = ?", *isMainGuardian)
	}

	if err := req.SetSearchBy(query, []string{"relationship", "note"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "relationship", "is_main_guardian", "is_emergency_contact"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateStudentGuardianByID(ctx context.Context, id uuid.UUID, data *ent.StudentGuardianUpdate) (*ent.StudentGuardian, error) {
	query := s.db.NewUpdate().
		Model(&ent.StudentGuardian{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.StudentID != nil {
		query.Set("student_id = ?", *data.StudentID)
	}
	if data.GuardianID != nil {
		query.Set("guardian_id = ?", *data.GuardianID)
	}
	if data.Relationship != nil {
		query.Set("relationship = ?", *data.Relationship)
	}
	if data.IsMainGuardian != nil {
		query.Set("is_main_guardian = ?", *data.IsMainGuardian)
	}
	if data.CanPickup != nil {
		query.Set("can_pickup = ?", *data.CanPickup)
	}
	if data.IsEmergencyContact != nil {
		query.Set("is_emergency_contact = ?", *data.IsEmergencyContact)
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

	return s.GetStudentGuardianByID(ctx, id)
}

func (s *Service) SoftDeleteStudentGuardianByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.StudentGuardian{}).
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

		_, err = tx.NewDelete().Model(&ent.StudentGuardian{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}

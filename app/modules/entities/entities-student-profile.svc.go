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

var _ entitiesinf.StudentProfileEntity = (*Service)(nil)

func (s *Service) CreateStudentProfile(ctx context.Context, data *ent.StudentProfile) (*ent.StudentProfile, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetStudentProfileByID(ctx context.Context, id uuid.UUID) (*ent.StudentProfile, error) {
	row := new(ent.StudentProfile)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListStudentProfiles(ctx context.Context, req *base.RequestPaginate, studentID *uuid.UUID) ([]*ent.StudentProfile, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.StudentProfile, 0)
	query := s.db.NewSelect().Model(&items)

	if studentID != nil {
		query.Where("student_id = ?", *studentID)
	}

	if err := req.SetSearchBy(query, []string{"nationality", "religion", "emergency_contact_name", "emergency_contact_phone"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "birth_date", "nationality", "religion"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateStudentProfileByID(ctx context.Context, id uuid.UUID, data *ent.StudentProfileUpdate) (*ent.StudentProfile, error) {
	query := s.db.NewUpdate().
		Model(&ent.StudentProfile{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.StudentID != nil {
		query.Set("student_id = ?", *data.StudentID)
	}
	if data.BirthDate != nil {
		query.Set("birth_date = ?", *data.BirthDate)
	}
	if data.Nationality != nil {
		query.Set("nationality = ?", *data.Nationality)
	}
	if data.Religion != nil {
		query.Set("religion = ?", *data.Religion)
	}
	if data.AddressCurrent != nil {
		query.Set("address_current = ?", *data.AddressCurrent)
	}
	if data.AddressRegistered != nil {
		query.Set("address_registered = ?", *data.AddressRegistered)
	}
	if data.EmergencyContactName != nil {
		query.Set("emergency_contact_name = ?", *data.EmergencyContactName)
	}
	if data.EmergencyContactPhone != nil {
		query.Set("emergency_contact_phone = ?", *data.EmergencyContactPhone)
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

	return s.GetStudentProfileByID(ctx, id)
}

func (s *Service) SoftDeleteStudentProfileByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.StudentProfile{}).
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

		_, err = tx.NewDelete().Model(&ent.StudentProfile{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}

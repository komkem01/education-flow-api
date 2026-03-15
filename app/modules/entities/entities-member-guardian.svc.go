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

var _ entitiesinf.MemberGuardianEntity = (*Service)(nil)

func (s *Service) CreateMemberGuardian(ctx context.Context, data *ent.MemberGuardian) (*ent.MemberGuardian, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetMemberGuardianByID(ctx context.Context, id uuid.UUID) (*ent.MemberGuardian, error) {
	row := new(ent.MemberGuardian)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListMemberGuardians(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, memberID *uuid.UUID) ([]*ent.MemberGuardian, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.MemberGuardian, 0)
	query := s.db.NewSelect().Model(&items)

	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}
	if schoolID != nil {
		query.Where("school_id = ?", *schoolID)
	}
	if memberID != nil {
		query.Where("member_id = ?", *memberID)
	}

	if err := req.SetSearchBy(query, []string{"first_name_th", "last_name_th", "first_name_en", "last_name_en", "citizen_id", "phone"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "first_name_th", "last_name_th", "is_active"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateMemberGuardianByID(ctx context.Context, id uuid.UUID, data *ent.MemberGuardianUpdate) (*ent.MemberGuardian, error) {
	query := s.db.NewUpdate().
		Model(&ent.MemberGuardian{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.MemberID != nil {
		query.Set("member_id = ?", *data.MemberID)
	}
	if data.SchoolID != nil {
		query.Set("school_id = ?", *data.SchoolID)
	}
	if data.GenderID != nil {
		query.Set("gender_id = ?", *data.GenderID)
	}
	if data.PrefixID != nil {
		query.Set("prefix_id = ?", *data.PrefixID)
	}
	if data.FirstNameTH != nil {
		query.Set("first_name_th = ?", *data.FirstNameTH)
	}
	if data.LastNameTH != nil {
		query.Set("last_name_th = ?", *data.LastNameTH)
	}
	if data.FirstNameEN != nil {
		query.Set("first_name_en = ?", *data.FirstNameEN)
	}
	if data.LastNameEN != nil {
		query.Set("last_name_en = ?", *data.LastNameEN)
	}
	if data.CitizenID != nil {
		query.Set("citizen_id = ?", *data.CitizenID)
	}
	if data.Phone != nil {
		query.Set("phone = ?", *data.Phone)
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

	return s.GetMemberGuardianByID(ctx, id)
}

func (s *Service) SoftDeleteMemberGuardianByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.MemberGuardian{}).
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

		_, err = tx.NewDelete().Model(&ent.MemberGuardian{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}

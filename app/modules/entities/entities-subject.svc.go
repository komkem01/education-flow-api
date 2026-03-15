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

var _ entitiesinf.SubjectEntity = (*Service)(nil)

func (s *Service) CreateSubject(ctx context.Context, data *ent.Subject) (*ent.Subject, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetSubjectByID(ctx context.Context, id uuid.UUID) (*ent.Subject, error) {
	row := new(ent.Subject)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListSubjects(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, subjectGroupID *uuid.UUID, isElective *bool) ([]*ent.Subject, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.Subject, 0)
	query := s.db.NewSelect().Model(&items)

	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}
	if schoolID != nil {
		query.Where("school_id = ?", *schoolID)
	}
	if subjectGroupID != nil {
		query.Where("subject_group_id = ?", *subjectGroupID)
	}
	if isElective != nil {
		query.Where("is_elective = ?", *isElective)
	}

	if err := req.SetSearchBy(query, []string{"code", "name_th", "name_en"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "code", "name_th", "credit", "is_elective", "is_active"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateSubjectByID(ctx context.Context, id uuid.UUID, data *ent.SubjectUpdate) (*ent.Subject, error) {
	query := s.db.NewUpdate().
		Model(&ent.Subject{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.SchoolID != nil {
		query.Set("school_id = ?", *data.SchoolID)
	}
	if data.SubjectGroupID != nil {
		query.Set("subject_group_id = ?", *data.SubjectGroupID)
	}
	if data.Code != nil {
		query.Set("code = ?", *data.Code)
	}
	if data.NameTH != nil {
		query.Set("name_th = ?", *data.NameTH)
	}
	if data.NameEN != nil {
		query.Set("name_en = ?", *data.NameEN)
	}
	if data.Credit != nil {
		query.Set("credit = ?", *data.Credit)
	}
	if data.HoursPerWeek != nil {
		query.Set("hours_per_week = ?", *data.HoursPerWeek)
	}
	if data.IsElective != nil {
		query.Set("is_elective = ?", *data.IsElective)
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

	return s.GetSubjectByID(ctx, id)
}

func (s *Service) SoftDeleteSubjectByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.Subject{}).
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

		_, err = tx.NewDelete().Model(&ent.Subject{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}

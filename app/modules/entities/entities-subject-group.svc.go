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

var _ entitiesinf.SubjectGroupEntity = (*Service)(nil)

func (s *Service) CreateSubjectGroup(ctx context.Context, data *ent.SubjectGroup) (*ent.SubjectGroup, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetSubjectGroupByID(ctx context.Context, id uuid.UUID) (*ent.SubjectGroup, error) {
	row := new(ent.SubjectGroup)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListSubjectGroups(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, headTeacherID *uuid.UUID) ([]*ent.SubjectGroup, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.SubjectGroup, 0)
	query := s.db.NewSelect().Model(&items)

	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}
	if schoolID != nil {
		query.Where("school_id = ?", *schoolID)
	}
	if headTeacherID != nil {
		query.Where("head_teacher_id = ?", *headTeacherID)
	}

	if err := req.SetSearchBy(query, []string{"code", "name_th", "name_en"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "code", "name_th", "is_active"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateSubjectGroupByID(ctx context.Context, id uuid.UUID, data *ent.SubjectGroupUpdate) (*ent.SubjectGroup, error) {
	query := s.db.NewUpdate().
		Model(&ent.SubjectGroup{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.SchoolID != nil {
		query.Set("school_id = ?", *data.SchoolID)
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
	if data.HeadTeacherID != nil {
		query.Set("head_teacher_id = ?", *data.HeadTeacherID)
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

	return s.GetSubjectGroupByID(ctx, id)
}

func (s *Service) SoftDeleteSubjectGroupByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.SubjectGroup{}).
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

		_, err = tx.NewDelete().Model(&ent.SubjectGroup{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}

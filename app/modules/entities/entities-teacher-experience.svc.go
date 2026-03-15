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

var _ entitiesinf.TeacherExperienceEntity = (*Service)(nil)

func (s *Service) CreateTeacherExperience(ctx context.Context, data *ent.TeacherExperience) (*ent.TeacherExperience, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetTeacherExperienceByID(ctx context.Context, id uuid.UUID) (*ent.TeacherExperience, error) {
	row := new(ent.TeacherExperience)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListTeacherExperiences(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, isCurrent *bool, isActive *bool) ([]*ent.TeacherExperience, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.TeacherExperience, 0)
	query := s.db.NewSelect().Model(&items)

	if teacherID != nil {
		query.Where("teacher_id = ?", *teacherID)
	}
	if isCurrent != nil {
		query.Where("is_current = ?", *isCurrent)
	}
	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}

	if err := req.SetSearchBy(query, []string{"school_name", "position", "department_name", "responsibilities", "achievements"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("sort_order ASC").Order("start_date DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "start_date", "sort_order", "school_name", "position", "is_current", "is_active"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateTeacherExperienceByID(ctx context.Context, id uuid.UUID, data *ent.TeacherExperienceUpdate) (*ent.TeacherExperience, error) {
	query := s.db.NewUpdate().
		Model(&ent.TeacherExperience{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.TeacherID != nil {
		query.Set("teacher_id = ?", *data.TeacherID)
	}
	if data.SchoolName != nil {
		query.Set("school_name = ?", *data.SchoolName)
	}
	if data.Position != nil {
		query.Set("position = ?", *data.Position)
	}
	if data.DepartmentName != nil {
		query.Set("department_name = ?", *data.DepartmentName)
	}
	if data.StartDate != nil {
		query.Set("start_date = ?", *data.StartDate)
	}
	if data.EndDate != nil {
		query.Set("end_date = ?", *data.EndDate)
	}
	if data.IsCurrent != nil {
		query.Set("is_current = ?", *data.IsCurrent)
	}
	if data.Responsibilities != nil {
		query.Set("responsibilities = ?", *data.Responsibilities)
	}
	if data.Achievements != nil {
		query.Set("achievements = ?", *data.Achievements)
	}
	if data.SortOrder != nil {
		query.Set("sort_order = ?", *data.SortOrder)
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

	return s.GetTeacherExperienceByID(ctx, id)
}

func (s *Service) SoftDeleteTeacherExperienceByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.TeacherExperience{}).
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

		_, err = tx.NewDelete().Model(&ent.TeacherExperience{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}

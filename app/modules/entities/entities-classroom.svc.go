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

var _ entitiesinf.ClassroomEntity = (*Service)(nil)

func (s *Service) CreateClassroom(ctx context.Context, data *ent.Classroom) (*ent.Classroom, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetClassroomByID(ctx context.Context, id uuid.UUID) (*ent.Classroom, error) {
	row := new(ent.Classroom)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListClassrooms(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, academicYearID *uuid.UUID, homeroomTeacherID *uuid.UUID) ([]*ent.Classroom, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.Classroom, 0)
	query := s.db.NewSelect().Model(&items)

	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}
	if schoolID != nil {
		query.Where("school_id = ?", *schoolID)
	}
	if academicYearID != nil {
		query.Where("academic_year_id = ?", *academicYearID)
	}
	if homeroomTeacherID != nil {
		query.Where("homeroom_teacher_id = ?", *homeroomTeacherID)
	}

	if err := req.SetSearchBy(query, []string{"name", "level", "room_no"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "name", "level", "room_no", "capacity", "is_active"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateClassroomByID(ctx context.Context, id uuid.UUID, data *ent.ClassroomUpdate) (*ent.Classroom, error) {
	query := s.db.NewUpdate().
		Model(&ent.Classroom{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.SchoolID != nil {
		query.Set("school_id = ?", *data.SchoolID)
	}
	if data.AcademicYearID != nil {
		query.Set("academic_year_id = ?", *data.AcademicYearID)
	}
	if data.Level != nil {
		query.Set("level = ?", *data.Level)
	}
	if data.RoomNo != nil {
		query.Set("room_no = ?", *data.RoomNo)
	}
	if data.Name != nil {
		query.Set("name = ?", *data.Name)
	}
	if data.HomeroomTeacherID != nil {
		query.Set("homeroom_teacher_id = ?", *data.HomeroomTeacherID)
	}
	if data.Capacity != nil {
		query.Set("capacity = ?", *data.Capacity)
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

	return s.GetClassroomByID(ctx, id)
}

func (s *Service) SoftDeleteClassroomByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.Classroom{}).
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

		_, err = tx.NewDelete().Model(&ent.Classroom{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}

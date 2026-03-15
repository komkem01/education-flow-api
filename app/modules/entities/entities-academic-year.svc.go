package entities

import (
	"context"
	"database/sql"
	"time"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.AcademicYearEntity = (*Service)(nil)

func (s *Service) CreateAcademicYear(ctx context.Context, schoolID uuid.UUID, year string, startDate time.Time, endDate time.Time, isActive bool) (*ent.AcademicYear, error) {
	academicYear := &ent.AcademicYear{
		SchoolID:  schoolID,
		Year:      year,
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  isActive,
	}

	if _, err := s.db.NewInsert().Model(academicYear).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return academicYear, nil
}

func (s *Service) GetAcademicYearByID(ctx context.Context, id uuid.UUID) (*ent.AcademicYear, error) {
	academicYear := new(ent.AcademicYear)
	if err := s.db.NewSelect().Model(academicYear).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return academicYear, nil
}

func (s *Service) ListAcademicYears(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID) ([]*ent.AcademicYear, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.AcademicYear, 0)
	query := s.db.NewSelect().Model(&items)

	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}

	if schoolID != nil {
		query.Where("school_id = ?", *schoolID)
	}

	if err := req.SetSearchBy(query, []string{"year"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}

	if err := req.SetSortOrder(query, []string{"created_at", "year", "start_date", "end_date", "is_active"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)

	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{
		Page:  req.GetPage(),
		Size:  req.GetSize(),
		Total: int64(total),
	}, nil
}

func (s *Service) UpdateAcademicYearByID(ctx context.Context, id uuid.UUID, schoolID *uuid.UUID, year *string, startDate *time.Time, endDate *time.Time, isActive *bool) (*ent.AcademicYear, error) {
	query := s.db.NewUpdate().
		Model(&ent.AcademicYear{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if schoolID != nil {
		query.Set("school_id = ?", *schoolID)
	}
	if year != nil {
		query.Set("year = ?", *year)
	}
	if startDate != nil {
		query.Set("start_date = ?", *startDate)
	}
	if endDate != nil {
		query.Set("end_date = ?", *endDate)
	}
	if isActive != nil {
		query.Set("is_active = ?", *isActive)
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

	return s.GetAcademicYearByID(ctx, id)
}

func (s *Service) SoftDeleteAcademicYearByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.AcademicYear{}).
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

		_, err = tx.NewDelete().Model(&ent.AcademicYear{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}

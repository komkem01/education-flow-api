package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.EnrollmentStatusHistoryEntity = (*Service)(nil)

func (s *Service) CreateEnrollmentStatusHistory(ctx context.Context, data *ent.EnrollmentStatusHistory) (*ent.EnrollmentStatusHistory, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetEnrollmentStatusHistoryByID(ctx context.Context, id uuid.UUID) (*ent.EnrollmentStatusHistory, error) {
	row := new(ent.EnrollmentStatusHistory)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListEnrollmentStatusHistories(ctx context.Context, req *base.RequestPaginate, enrollmentID *uuid.UUID, toStatus *ent.StudentEnrollmentStatus, changedBy *uuid.UUID) ([]*ent.EnrollmentStatusHistory, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.EnrollmentStatusHistory, 0)
	query := s.db.NewSelect().Model(&items)

	if enrollmentID != nil {
		query.Where("enrollment_id = ?", *enrollmentID)
	}
	if toStatus != nil {
		query.Where("to_status = ?", *toStatus)
	}
	if changedBy != nil {
		query.Where("changed_by = ?", *changedBy)
	}

	if err := req.SetSearchBy(query, []string{"reason"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("changed_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"changed_at", "created_at", "to_status"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateEnrollmentStatusHistoryByID(ctx context.Context, id uuid.UUID, data *ent.EnrollmentStatusHistoryUpdate) (*ent.EnrollmentStatusHistory, error) {
	query := s.db.NewUpdate().
		Model(&ent.EnrollmentStatusHistory{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.EnrollmentID != nil {
		query.Set("enrollment_id = ?", *data.EnrollmentID)
	}
	if data.FromStatus != nil {
		query.Set("from_status = ?", *data.FromStatus)
	}
	if data.ToStatus != nil {
		query.Set("to_status = ?", *data.ToStatus)
	}
	if data.ChangedAt != nil {
		query.Set("changed_at = ?", *data.ChangedAt)
	}
	if data.ChangedBy != nil {
		query.Set("changed_by = ?", *data.ChangedBy)
	}
	if data.Reason != nil {
		query.Set("reason = ?", *data.Reason)
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

	return s.GetEnrollmentStatusHistoryByID(ctx, id)
}

func (s *Service) SoftDeleteEnrollmentStatusHistoryByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.EnrollmentStatusHistory{}).Where("id = ?", id).Exec(ctx)
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

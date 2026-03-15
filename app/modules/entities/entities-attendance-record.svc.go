package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.AttendanceRecordEntity = (*Service)(nil)

func (s *Service) CreateAttendanceRecord(ctx context.Context, data *ent.AttendanceRecord) (*ent.AttendanceRecord, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetAttendanceRecordByID(ctx context.Context, id uuid.UUID) (*ent.AttendanceRecord, error) {
	row := new(ent.AttendanceRecord)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) GetAttendanceRecordBySessionAndEnrollment(ctx context.Context, sessionID uuid.UUID, enrollmentID uuid.UUID) (*ent.AttendanceRecord, error) {
	row := new(ent.AttendanceRecord)
	if err := s.db.NewSelect().Model(row).Where("session_id = ?", sessionID).Where("enrollment_id = ?", enrollmentID).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListAttendanceRecords(ctx context.Context, req *base.RequestPaginate, sessionID *uuid.UUID, enrollmentID *uuid.UUID, status *ent.AttendanceStatus, source *ent.AttendanceSource, markedBy *uuid.UUID) ([]*ent.AttendanceRecord, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.AttendanceRecord, 0)
	query := s.db.NewSelect().Model(&items)

	if sessionID != nil {
		query.Where("session_id = ?", *sessionID)
	}
	if enrollmentID != nil {
		query.Where("enrollment_id = ?", *enrollmentID)
	}
	if status != nil {
		query.Where("status = ?", *status)
	}
	if source != nil {
		query.Where("source = ?", *source)
	}
	if markedBy != nil {
		query.Where("marked_by = ?", *markedBy)
	}

	if err := req.SetSearchBy(query, []string{"status", "source", "remark"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("marked_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"marked_at", "status", "source", "created_at"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateAttendanceRecordByID(ctx context.Context, id uuid.UUID, data *ent.AttendanceRecordUpdate) (*ent.AttendanceRecord, error) {
	query := s.db.NewUpdate().
		Model(&ent.AttendanceRecord{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.SessionID != nil {
		query.Set("session_id = ?", *data.SessionID)
	}
	if data.EnrollmentID != nil {
		query.Set("enrollment_id = ?", *data.EnrollmentID)
	}
	if data.Status != nil {
		query.Set("status = ?", *data.Status)
	}
	if data.Source != nil {
		query.Set("source = ?", *data.Source)
	}
	if data.MarkedAt != nil {
		query.Set("marked_at = ?", *data.MarkedAt)
	}
	if data.Remark != nil {
		query.Set("remark = ?", *data.Remark)
	}
	if data.MarkedBy != nil {
		query.Set("marked_by = ?", *data.MarkedBy)
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

	return s.GetAttendanceRecordByID(ctx, id)
}

func (s *Service) SoftDeleteAttendanceRecordByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.AttendanceRecord{}).Where("id = ?", id).Exec(ctx)
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

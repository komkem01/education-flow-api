package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.AttendanceRecordLogEntity = (*Service)(nil)

func (s *Service) CreateAttendanceRecordLog(ctx context.Context, data *ent.AttendanceRecordLog) (*ent.AttendanceRecordLog, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetAttendanceRecordLogByID(ctx context.Context, id uuid.UUID) (*ent.AttendanceRecordLog, error) {
	row := new(ent.AttendanceRecordLog)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListAttendanceRecordLogs(ctx context.Context, req *base.RequestPaginate, recordID *uuid.UUID, changedBy *uuid.UUID, newStatus *ent.AttendanceStatus) ([]*ent.AttendanceRecordLog, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.AttendanceRecordLog, 0)
	query := s.db.NewSelect().Model(&items)

	if recordID != nil {
		query.Where("record_id = ?", *recordID)
	}
	if changedBy != nil {
		query.Where("changed_by = ?", *changedBy)
	}
	if newStatus != nil {
		query.Where("new_status = ?", *newStatus)
	}

	if err := req.SetSearchBy(query, []string{"reason", "new_status"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("changed_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"changed_at", "new_status", "created_at"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateAttendanceRecordLogByID(ctx context.Context, id uuid.UUID, data *ent.AttendanceRecordLogUpdate) (*ent.AttendanceRecordLog, error) {
	query := s.db.NewUpdate().
		Model(&ent.AttendanceRecordLog{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.RecordID != nil {
		query.Set("record_id = ?", *data.RecordID)
	}
	if data.OldStatus != nil {
		query.Set("old_status = ?", *data.OldStatus)
	}
	if data.NewStatus != nil {
		query.Set("new_status = ?", *data.NewStatus)
	}
	if data.ChangedBy != nil {
		query.Set("changed_by = ?", *data.ChangedBy)
	}
	if data.ChangedAt != nil {
		query.Set("changed_at = ?", *data.ChangedAt)
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

	return s.GetAttendanceRecordLogByID(ctx, id)
}

func (s *Service) SoftDeleteAttendanceRecordLogByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.AttendanceRecordLog{}).Where("id = ?", id).Exec(ctx)
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

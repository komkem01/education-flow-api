package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type AttendanceRecordLogEntity interface {
	CreateAttendanceRecordLog(ctx context.Context, data *ent.AttendanceRecordLog) (*ent.AttendanceRecordLog, error)
	GetAttendanceRecordLogByID(ctx context.Context, id uuid.UUID) (*ent.AttendanceRecordLog, error)
	ListAttendanceRecordLogs(ctx context.Context, req *base.RequestPaginate, recordID *uuid.UUID, changedBy *uuid.UUID, newStatus *ent.AttendanceStatus) ([]*ent.AttendanceRecordLog, *base.ResponsePaginate, error)
	UpdateAttendanceRecordLogByID(ctx context.Context, id uuid.UUID, data *ent.AttendanceRecordLogUpdate) (*ent.AttendanceRecordLog, error)
	SoftDeleteAttendanceRecordLogByID(ctx context.Context, id uuid.UUID) error
}

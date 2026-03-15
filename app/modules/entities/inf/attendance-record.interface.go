package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type AttendanceRecordEntity interface {
	CreateAttendanceRecord(ctx context.Context, data *ent.AttendanceRecord) (*ent.AttendanceRecord, error)
	GetAttendanceRecordByID(ctx context.Context, id uuid.UUID) (*ent.AttendanceRecord, error)
	GetAttendanceRecordBySessionAndEnrollment(ctx context.Context, sessionID uuid.UUID, enrollmentID uuid.UUID) (*ent.AttendanceRecord, error)
	ListAttendanceRecords(ctx context.Context, req *base.RequestPaginate, sessionID *uuid.UUID, enrollmentID *uuid.UUID, status *ent.AttendanceStatus, source *ent.AttendanceSource, markedBy *uuid.UUID) ([]*ent.AttendanceRecord, *base.ResponsePaginate, error)
	UpdateAttendanceRecordByID(ctx context.Context, id uuid.UUID, data *ent.AttendanceRecordUpdate) (*ent.AttendanceRecord, error)
	SoftDeleteAttendanceRecordByID(ctx context.Context, id uuid.UUID) error
}

package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type AttendanceSessionEntity interface {
	CreateAttendanceSession(ctx context.Context, data *ent.AttendanceSession) (*ent.AttendanceSession, error)
	GetAttendanceSessionByID(ctx context.Context, id uuid.UUID) (*ent.AttendanceSession, error)
	ListAttendanceSessions(ctx context.Context, req *base.RequestPaginate, schoolID *uuid.UUID, academicYearID *uuid.UUID, classroomID *uuid.UUID, subjectID *uuid.UUID, teacherID *uuid.UUID, mode *ent.AttendanceMode, sessionDateFrom *string, sessionDateTo *string) ([]*ent.AttendanceSession, *base.ResponsePaginate, error)
	UpdateAttendanceSessionByID(ctx context.Context, id uuid.UUID, data *ent.AttendanceSessionUpdate) (*ent.AttendanceSession, error)
	SoftDeleteAttendanceSessionByID(ctx context.Context, id uuid.UUID) error
}

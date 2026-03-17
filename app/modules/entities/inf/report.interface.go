package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"

	"github.com/google/uuid"
)

type ReportEntity interface {
	CountMembersByRole(ctx context.Context, schoolID uuid.UUID, role ent.MemberRole, isActive *bool) (int64, error)
	CountSubjectsBySchool(ctx context.Context, schoolID uuid.UUID) (int64, error)
	CountAttendanceByStatus(ctx context.Context, schoolID uuid.UUID, academicYearID *uuid.UUID, status *ent.AttendanceStatus) (int64, error)
	ListAcademicYearsBySchool(ctx context.Context, schoolID uuid.UUID) ([]*ent.AcademicYear, error)
}

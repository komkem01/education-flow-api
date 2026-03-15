package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type EnrollmentStatusHistoryEntity interface {
	CreateEnrollmentStatusHistory(ctx context.Context, data *ent.EnrollmentStatusHistory) (*ent.EnrollmentStatusHistory, error)
	GetEnrollmentStatusHistoryByID(ctx context.Context, id uuid.UUID) (*ent.EnrollmentStatusHistory, error)
	ListEnrollmentStatusHistories(ctx context.Context, req *base.RequestPaginate, enrollmentID *uuid.UUID, toStatus *ent.StudentEnrollmentStatus, changedBy *uuid.UUID) ([]*ent.EnrollmentStatusHistory, *base.ResponsePaginate, error)
	UpdateEnrollmentStatusHistoryByID(ctx context.Context, id uuid.UUID, data *ent.EnrollmentStatusHistoryUpdate) (*ent.EnrollmentStatusHistory, error)
	SoftDeleteEnrollmentStatusHistoryByID(ctx context.Context, id uuid.UUID) error
}

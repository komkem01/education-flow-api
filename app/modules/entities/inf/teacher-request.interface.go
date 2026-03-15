package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type TeacherRequestEntity interface {
	CreateTeacherRequest(ctx context.Context, data *ent.TeacherRequest) (*ent.TeacherRequest, error)
	GetTeacherRequestByID(ctx context.Context, id uuid.UUID) (*ent.TeacherRequest, error)
	ListTeacherRequests(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, requestType *ent.TeacherRequestType, status *ent.TeacherRequestStatus) ([]*ent.TeacherRequest, *base.ResponsePaginate, error)
	UpdateTeacherRequestByID(ctx context.Context, id uuid.UUID, data *ent.TeacherRequestUpdate) (*ent.TeacherRequest, error)
	DeleteTeacherRequestByID(ctx context.Context, id uuid.UUID) error
}

package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type ApprovalRequestEntity interface {
	CreateApprovalRequest(ctx context.Context, data *ent.ApprovalRequest) (*ent.ApprovalRequest, error)
	GetApprovalRequestByID(ctx context.Context, id uuid.UUID) (*ent.ApprovalRequest, error)
	ListApprovalRequests(ctx context.Context, req *base.RequestPaginate, requestedBy *uuid.UUID, requestedByRole *ent.ApprovalActorRole, status *ent.ApprovalRequestStatus, requestType *string) ([]*ent.ApprovalRequest, *base.ResponsePaginate, error)
	UpdateApprovalRequestByID(ctx context.Context, id uuid.UUID, data *ent.ApprovalRequestUpdate) (*ent.ApprovalRequest, error)
	DeleteApprovalRequestByID(ctx context.Context, id uuid.UUID) error
}

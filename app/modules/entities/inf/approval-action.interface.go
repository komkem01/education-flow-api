package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type ApprovalActionEntity interface {
	CreateApprovalAction(ctx context.Context, data *ent.ApprovalAction) (*ent.ApprovalAction, error)
	GetApprovalActionByIdempotencyKey(ctx context.Context, requestID uuid.UUID, action ent.ApprovalActionType, idempotencyKey string) (*ent.ApprovalAction, error)
	ListApprovalActions(ctx context.Context, req *base.RequestPaginate, requestID *uuid.UUID, actedBy *uuid.UUID, action *ent.ApprovalActionType) ([]*ent.ApprovalAction, *base.ResponsePaginate, error)
}

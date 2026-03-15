package approvals

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Reject(ctx context.Context, id uuid.UUID, actorID string, actorRole string, idempotencyKey *string, comment *string, metadata map[string]any) (*ent.ApprovalRequest, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "approvals.service.reject")
	defer span.End()

	parsedActorID, parsedRole, err := s.parseActor(actorID, actorRole)
	if err != nil {
		return nil, err
	}
	if parsedRole != ent.ApprovalActorRoleAdmin {
		return nil, fmt.Errorf("%w", ErrApprovalRequestUnauthorized)
	}

	idempotentReq, err := s.findIdempotentRequest(ctx, id, ent.ApprovalActionTypeReject, idempotencyKey)
	if err != nil {
		return nil, err
	}
	if idempotentReq != nil {
		return idempotentReq, nil
	}

	item, err := s.loadRequest(ctx, id)
	if err != nil {
		return nil, err
	}
	if item.CurrentStatus != ent.ApprovalRequestStatusPending {
		return nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}

	now := time.Now()
	updated, err := s.db.UpdateApprovalRequestByID(ctx, id, &ent.ApprovalRequestUpdate{
		CurrentStatus: ptrApprovalStatus(ent.ApprovalRequestStatusRejected),
		ResolvedAt:    &now,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	if err := s.appendActionWithIdempotency(ctx, id, ent.ApprovalActionTypeReject, parsedActorID, parsedRole, idempotencyKey, comment, metadata); err != nil {
		return nil, err
	}

	return updated, nil
}

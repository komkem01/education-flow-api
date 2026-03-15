package approvals

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Cancel(ctx context.Context, id uuid.UUID, actorID string, actorRole string, idempotencyKey *string, comment *string, metadata map[string]any) (*ent.ApprovalRequest, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "approvals.service.cancel")
	defer span.End()

	parsedActorID, parsedRole, err := s.parseActor(actorID, actorRole)
	if err != nil {
		return nil, err
	}

	idempotentReq, err := s.findIdempotentRequest(ctx, id, ent.ApprovalActionTypeCancel, idempotencyKey)
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
	if isTerminalStatus(item.CurrentStatus) {
		return nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}

	switch parsedRole {
	case ent.ApprovalActorRoleAdmin:
		// Management can cancel any non-terminal request.
	case ent.ApprovalActorRoleTeacher:
		if parsedActorID != item.RequestedBy {
			return nil, fmt.Errorf("%w", ErrApprovalRequestUnauthorized)
		}
		if item.CurrentStatus != ent.ApprovalRequestStatusDraft && item.CurrentStatus != ent.ApprovalRequestStatusPending {
			return nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
	default:
		return nil, fmt.Errorf("%w", ErrApprovalRequestUnauthorized)
	}

	now := time.Now()
	updated, err := s.db.UpdateApprovalRequestByID(ctx, id, &ent.ApprovalRequestUpdate{
		CurrentStatus: ptrApprovalStatus(ent.ApprovalRequestStatusCancelled),
		ResolvedAt:    &now,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	if err := s.appendActionWithIdempotency(ctx, id, ent.ApprovalActionTypeCancel, parsedActorID, parsedRole, idempotencyKey, comment, metadata); err != nil {
		return nil, err
	}

	return updated, nil
}

package approvals

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Submit(ctx context.Context, id uuid.UUID, actorID string, actorRole string, comment *string) (*ent.ApprovalRequest, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "approvals.service.submit")
	defer span.End()

	parsedActorID, parsedRole, err := s.parseActor(actorID, actorRole)
	if err != nil {
		return nil, err
	}

	item, err := s.loadRequest(ctx, id)
	if err != nil {
		return nil, err
	}

	if item.CurrentStatus != ent.ApprovalRequestStatusDraft {
		return nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}
	if item.RequestedBy != parsedActorID {
		return nil, fmt.Errorf("%w", ErrApprovalRequestUnauthorized)
	}
	if parsedRole != item.RequestedByRole {
		return nil, fmt.Errorf("%w", ErrApprovalRequestUnauthorized)
	}

	now := time.Now()
	updated, err := s.db.UpdateApprovalRequestByID(ctx, id, &ent.ApprovalRequestUpdate{
		CurrentStatus: ptrApprovalStatus(ent.ApprovalRequestStatusPending),
		SubmittedAt:   &now,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	if err := s.appendAction(ctx, id, ent.ApprovalActionTypeSubmit, parsedActorID, parsedRole, comment, nil); err != nil {
		return nil, err
	}

	return updated, nil
}

func ptrApprovalStatus(v ent.ApprovalRequestStatus) *ent.ApprovalRequestStatus {
	return &v
}

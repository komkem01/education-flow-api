package approvals

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Approve(ctx context.Context, id uuid.UUID, actorID string, actorRole string, idempotencyKey *string, comment *string, metadata map[string]any) (*ent.ApprovalRequest, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "approvals.service.approve")
	defer span.End()

	parsedActorID, parsedRole, err := s.parseActor(actorID, actorRole)
	if err != nil {
		return nil, err
	}
	if parsedRole != ent.ApprovalActorRoleAdmin {
		return nil, fmt.Errorf("%w", ErrApprovalRequestUnauthorized)
	}

	idempotentReq, err := s.findIdempotentRequest(ctx, id, ent.ApprovalActionTypeApprove, idempotencyKey)
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

	if item.RequestType == "student_registration_case" || item.SubjectType == "student_registration_case" {
		if s.studentRegDB == nil || item.SubjectID == nil {
			return nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}

		if _, err := s.studentRegDB.ApproveAndApplyStudentRegistrationCase(
			ctx,
			*item.SubjectID,
			parsedActorID,
			parsedRole,
			comment,
			idempotencyKey,
			metadata,
		); err != nil {
			return nil, normalizeServiceError(err)
		}

		return s.loadRequest(ctx, id)
	}

	if err := s.applyApprovalEffects(ctx, item); err != nil {
		return nil, err
	}

	now := time.Now()
	updated, err := s.db.UpdateApprovalRequestByID(ctx, id, &ent.ApprovalRequestUpdate{
		CurrentStatus: ptrApprovalStatus(ent.ApprovalRequestStatusApproved),
		ResolvedAt:    &now,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	if err := s.appendActionWithIdempotency(ctx, id, ent.ApprovalActionTypeApprove, parsedActorID, parsedRole, idempotencyKey, comment, metadata); err != nil {
		return nil, err
	}

	return updated, nil
}

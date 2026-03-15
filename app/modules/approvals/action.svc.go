package approvals

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) parseActor(actorID string, actorRole string) (uuid.UUID, ent.ApprovalActorRole, error) {
	parsedID, err := uuid.Parse(actorID)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}
	parsedRole, ok := parseApprovalActorRole(actorRole)
	if !ok {
		return uuid.Nil, "", fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}
	return parsedID, parsedRole, nil
}

func (s *Service) appendAction(ctx context.Context, requestID uuid.UUID, action ent.ApprovalActionType, actorID uuid.UUID, actorRole ent.ApprovalActorRole, comment *string, metadata map[string]any) error {
	return s.appendActionWithIdempotency(ctx, requestID, action, actorID, actorRole, nil, comment, metadata)
}

func (s *Service) appendActionWithIdempotency(ctx context.Context, requestID uuid.UUID, action ent.ApprovalActionType, actorID uuid.UUID, actorRole ent.ApprovalActorRole, idempotencyKey *string, comment *string, metadata map[string]any) error {
	if metadata == nil {
		metadata = map[string]any{}
	}
	_, err := s.actionDB.CreateApprovalAction(ctx, &ent.ApprovalAction{
		RequestID:      requestID,
		Action:         action,
		IdempotencyKey: idempotencyKey,
		ActedBy:        actorID,
		ActedByRole:    actorRole,
		Comment:        comment,
		Metadata:       metadata,
		CreatedAt:      time.Now(),
	})
	if err != nil {
		return normalizeServiceError(err)
	}
	return nil
}

func (s *Service) findIdempotentRequest(ctx context.Context, requestID uuid.UUID, action ent.ApprovalActionType, idempotencyKey *string) (*ent.ApprovalRequest, error) {
	if idempotencyKey == nil || *idempotencyKey == "" {
		return nil, nil
	}

	_, err := s.actionDB.GetApprovalActionByIdempotencyKey(ctx, requestID, action, *idempotencyKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, normalizeServiceError(err)
	}

	return s.loadRequest(ctx, requestID)
}

func (s *Service) loadRequest(ctx context.Context, id uuid.UUID) (*ent.ApprovalRequest, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "approvals.service.load_request")
	defer span.End()

	item, err := s.db.GetApprovalRequestByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

package approvals

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, requestType string, subjectType string, subjectID *string, requestedBy string, requestedByRole string, payload map[string]any) (*ent.ApprovalRequest, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "approvals.service.create")
	defer span.End()

	if requestType == "" || subjectType == "" || requestedBy == "" {
		return nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}

	requestedByID, err := uuid.Parse(requestedBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}

	requestedRole, ok := parseApprovalActorRole(requestedByRole)
	if !ok {
		return nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}

	var subjectUUID *uuid.UUID
	if subjectID != nil && *subjectID != "" {
		parsed, err := uuid.Parse(*subjectID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		subjectUUID = &parsed
	}

	if payload == nil {
		payload = map[string]any{}
	}

	item, err := s.db.CreateApprovalRequest(ctx, &ent.ApprovalRequest{
		RequestType:     requestType,
		SubjectType:     subjectType,
		SubjectID:       subjectUUID,
		RequestedBy:     requestedByID,
		RequestedByRole: requestedRole,
		Payload:         payload,
		CurrentStatus:   ent.ApprovalRequestStatusDraft,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

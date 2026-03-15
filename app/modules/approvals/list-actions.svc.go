package approvals

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) ListActions(ctx context.Context, req *base.RequestPaginate, requestID uuid.UUID, actedBy *uuid.UUID, action *string) ([]*ent.ApprovalAction, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "approvals.service.list_actions")
	defer span.End()

	var actionVal *ent.ApprovalActionType
	if action != nil {
		parsed, ok := parseApprovalActionType(*action)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		actionVal = &parsed
	}

	items, page, err := s.actionDB.ListApprovalActions(ctx, req, &requestID, actedBy, actionVal)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrApprovalRequestConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

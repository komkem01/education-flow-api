package approvals

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, requestedBy *uuid.UUID, requestedByRole *string, status *string, requestType *string) ([]*ent.ApprovalRequest, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "approvals.service.list")
	defer span.End()

	var roleVal *ent.ApprovalActorRole
	if requestedByRole != nil {
		parsed, ok := parseApprovalActorRole(*requestedByRole)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		roleVal = &parsed
	}

	var statusVal *ent.ApprovalRequestStatus
	if status != nil {
		parsed, ok := parseApprovalStatus(*status)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
		}
		statusVal = &parsed
	}

	items, page, err := s.db.ListApprovalRequests(ctx, req, requestedBy, roleVal, statusVal, requestType)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrApprovalRequestConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

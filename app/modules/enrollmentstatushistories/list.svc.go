package enrollmentstatushistories

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, enrollmentID *uuid.UUID, toStatus *string, changedBy *uuid.UUID) ([]*ent.EnrollmentStatusHistory, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "enrollmentstatushistories.service.list")
	defer span.End()

	var toStatusVal *ent.StudentEnrollmentStatus
	if toStatus != nil {
		parsed, ok := parseStudentEnrollmentStatus(*toStatus)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrEnrollmentStatusHistoryConditionFail)
		}
		toStatusVal = &parsed
	}

	items, page, err := s.db.ListEnrollmentStatusHistories(ctx, req, enrollmentID, toStatusVal, changedBy)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrEnrollmentStatusHistoryConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

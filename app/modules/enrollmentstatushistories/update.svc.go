package enrollmentstatushistories

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.EnrollmentStatusHistory, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "enrollmentstatushistories.service.update")
	defer span.End()

	if req.EnrollmentID == nil && req.FromStatus == nil && req.ToStatus == nil && req.ChangedAt == nil && req.ChangedBy == nil && req.Reason == nil {
		return nil, fmt.Errorf("%w", ErrEnrollmentStatusHistoryConditionFail)
	}

	payload := &ent.EnrollmentStatusHistoryUpdate{Reason: req.Reason}

	var err error
	payload.EnrollmentID, err = parseOptionalUUID(req.EnrollmentID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrEnrollmentStatusHistoryConditionFail)
	}
	payload.ChangedBy, err = parseOptionalUUID(req.ChangedBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrEnrollmentStatusHistoryConditionFail)
	}
	payload.ChangedAt, err = parseOptionalDateTime(req.ChangedAt)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrEnrollmentStatusHistoryConditionFail)
	}

	if req.FromStatus != nil {
		if *req.FromStatus == "" {
			payload.FromStatus = nil
		} else {
			parsed, ok := parseStudentEnrollmentStatus(*req.FromStatus)
			if !ok {
				return nil, fmt.Errorf("%w", ErrEnrollmentStatusHistoryConditionFail)
			}
			payload.FromStatus = &parsed
		}
	}
	if req.ToStatus != nil {
		parsed, ok := parseStudentEnrollmentStatus(*req.ToStatus)
		if !ok {
			return nil, fmt.Errorf("%w", ErrEnrollmentStatusHistoryConditionFail)
		}
		payload.ToStatus = &parsed
	}

	item, err := s.db.UpdateEnrollmentStatusHistoryByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

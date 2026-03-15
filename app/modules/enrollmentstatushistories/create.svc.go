package enrollmentstatushistories

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, enrollmentID uuid.UUID, fromStatus *string, toStatus string, changedAt *string, changedBy *string, reason *string) (*ent.EnrollmentStatusHistory, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "enrollmentstatushistories.service.create")
	defer span.End()

	toStatusVal, ok := parseStudentEnrollmentStatus(toStatus)
	if !ok {
		return nil, fmt.Errorf("%w", ErrEnrollmentStatusHistoryConditionFail)
	}

	var fromStatusVal *ent.StudentEnrollmentStatus
	if fromStatus != nil && *fromStatus != "" {
		parsed, ok := parseStudentEnrollmentStatus(*fromStatus)
		if !ok {
			return nil, fmt.Errorf("%w", ErrEnrollmentStatusHistoryConditionFail)
		}
		fromStatusVal = &parsed
	}

	changedAtVal := time.Now()
	if changedAt != nil && *changedAt != "" {
		parsed, err := time.Parse(time.RFC3339, *changedAt)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrEnrollmentStatusHistoryConditionFail)
		}
		changedAtVal = parsed
	}

	changedByVal, err := parseOptionalUUID(changedBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrEnrollmentStatusHistoryConditionFail)
	}

	item, err := s.db.CreateEnrollmentStatusHistory(ctx, &ent.EnrollmentStatusHistory{
		EnrollmentID: enrollmentID,
		FromStatus:   fromStatusVal,
		ToStatus:     toStatusVal,
		ChangedAt:    changedAtVal,
		ChangedBy:    changedByVal,
		Reason:       reason,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

func parseOptionalUUID(v *string) (*uuid.UUID, error) {
	if v == nil || *v == "" {
		return nil, nil
	}
	parsed, err := uuid.Parse(*v)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseOptionalDateTime(v *string) (*time.Time, error) {
	if v == nil || *v == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, *v)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

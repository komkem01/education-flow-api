package teacherrequests

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, teacherID uuid.UUID, requestType string, requestData map[string]any, requestReason *string, status *string, approvedBy *uuid.UUID, approvedAt *string) (*ent.TeacherRequest, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacherrequests.service.create")
	defer span.End()

	parsedType, ok := parseTeacherRequestType(requestType)
	if !ok {
		return nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
	}

	parsedStatus := ent.TeacherRequestStatusPending
	if status != nil {
		v, ok := parseTeacherRequestStatus(*status)
		if !ok {
			return nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
		}
		parsedStatus = v
	}

	parsedApprovedAt, err := parseOptionalDateTime(approvedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTeacherRequestConditionFail, err)
	}

	item, err := s.db.CreateTeacherRequest(ctx, &ent.TeacherRequest{
		TeacherID:     teacherID,
		RequestType:   parsedType,
		RequestData:   requestData,
		RequestReason: requestReason,
		Status:        parsedStatus,
		ApprovedBy:    approvedBy,
		ApprovedAt:    parsedApprovedAt,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

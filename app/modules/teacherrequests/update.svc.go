package teacherrequests

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func parseOptionalDateTime(v *string) (*time.Time, error) {
	if v == nil || *v == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, *v)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.TeacherRequest, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacherrequests.service.update")
	defer span.End()

	if req.TeacherID == nil && req.RequestType == nil && req.RequestData == nil && req.RequestReason == nil && req.Status == nil && req.ApprovedBy == nil && req.ApprovedAt == nil {
		return nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
	}

	payload := &ent.TeacherRequestUpdate{
		RequestData:   req.RequestData,
		RequestReason: req.RequestReason,
	}

	if req.TeacherID != nil {
		parsed, err := uuid.Parse(*req.TeacherID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
		}
		payload.TeacherID = &parsed
	}
	if req.RequestType != nil {
		parsed, ok := parseTeacherRequestType(*req.RequestType)
		if !ok {
			return nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
		}
		payload.RequestType = &parsed
	}
	if req.Status != nil {
		parsed, ok := parseTeacherRequestStatus(*req.Status)
		if !ok {
			return nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
		}
		payload.Status = &parsed
	}
	if req.ApprovedBy != nil {
		parsed, err := uuid.Parse(*req.ApprovedBy)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
		}
		payload.ApprovedBy = &parsed
	}
	if req.ApprovedAt != nil {
		parsed, err := parseOptionalDateTime(req.ApprovedAt)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
		}
		payload.ApprovedAt = parsed
	}

	item, err := s.db.UpdateTeacherRequestByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

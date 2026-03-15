package attendancerecords

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.AttendanceRecord, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "attendancerecords.service.update")
	defer span.End()

	if req.SessionID == nil && req.EnrollmentID == nil && req.Status == nil && req.Source == nil && req.MarkedAt == nil && req.Remark == nil && req.MarkedBy == nil {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
	}

	payload := &ent.AttendanceRecordUpdate{Remark: req.Remark}

	var err error
	payload.SessionID, err = parseOptionalUUID(req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
	}
	payload.EnrollmentID, err = parseOptionalUUID(req.EnrollmentID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
	}
	payload.MarkedBy, err = parseOptionalUUID(req.MarkedBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
	}
	payload.MarkedAt, err = parseOptionalDateTime(req.MarkedAt)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
	}

	if req.Status != nil {
		parsed, ok := parseAttendanceStatus(*req.Status)
		if !ok {
			return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
		}
		payload.Status = &parsed
	}
	if req.Source != nil {
		parsed, ok := parseAttendanceSource(*req.Source)
		if !ok {
			return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
		}
		payload.Source = &parsed
	}

	item, err := s.db.UpdateAttendanceRecordByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

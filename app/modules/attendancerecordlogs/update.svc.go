package attendancerecordlogs

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.AttendanceRecordLog, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "attendancerecordlogs.service.update")
	defer span.End()

	if req.RecordID == nil && req.OldStatus == nil && req.NewStatus == nil && req.ChangedBy == nil && req.ChangedAt == nil && req.Reason == nil {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordLogConditionFail)
	}

	payload := &ent.AttendanceRecordLogUpdate{Reason: req.Reason}

	var err error
	payload.RecordID, err = parseOptionalUUID(req.RecordID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordLogConditionFail)
	}
	payload.ChangedBy, err = parseOptionalUUID(req.ChangedBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordLogConditionFail)
	}
	payload.ChangedAt, err = parseOptionalDateTime(req.ChangedAt)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordLogConditionFail)
	}

	if req.OldStatus != nil {
		if *req.OldStatus == "" {
			payload.OldStatus = nil
		} else {
			parsed, ok := parseAttendanceStatus(*req.OldStatus)
			if !ok {
				return nil, fmt.Errorf("%w", ErrAttendanceRecordLogConditionFail)
			}
			payload.OldStatus = &parsed
		}
	}
	if req.NewStatus != nil {
		parsed, ok := parseAttendanceStatus(*req.NewStatus)
		if !ok {
			return nil, fmt.Errorf("%w", ErrAttendanceRecordLogConditionFail)
		}
		payload.NewStatus = &parsed
	}

	item, err := s.db.UpdateAttendanceRecordLogByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

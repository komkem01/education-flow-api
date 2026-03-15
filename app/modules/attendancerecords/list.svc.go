package attendancerecords

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, sessionID *uuid.UUID, enrollmentID *uuid.UUID, status *string, source *string, markedBy *uuid.UUID) ([]*ent.AttendanceRecord, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "attendancerecords.service.list")
	defer span.End()

	var statusVal *ent.AttendanceStatus
	if status != nil {
		parsed, ok := parseAttendanceStatus(*status)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
		}
		statusVal = &parsed
	}

	var sourceVal *ent.AttendanceSource
	if source != nil {
		parsed, ok := parseAttendanceSource(*source)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
		}
		sourceVal = &parsed
	}

	items, page, err := s.db.ListAttendanceRecords(ctx, req, sessionID, enrollmentID, statusVal, sourceVal, markedBy)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrAttendanceRecordConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

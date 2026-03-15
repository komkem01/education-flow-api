package attendancerecordlogs

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, recordID *uuid.UUID, changedBy *uuid.UUID, newStatus *string) ([]*ent.AttendanceRecordLog, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "attendancerecordlogs.service.list")
	defer span.End()

	var statusVal *ent.AttendanceStatus
	if newStatus != nil {
		parsed, ok := parseAttendanceStatus(*newStatus)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrAttendanceRecordLogConditionFail)
		}
		statusVal = &parsed
	}

	items, page, err := s.db.ListAttendanceRecordLogs(ctx, req, recordID, changedBy, statusVal)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrAttendanceRecordLogConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

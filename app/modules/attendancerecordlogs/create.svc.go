package attendancerecordlogs

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, recordID uuid.UUID, oldStatus *string, newStatus string, changedBy *string, changedAt *string, reason *string) (*ent.AttendanceRecordLog, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "attendancerecordlogs.service.create")
	defer span.End()

	newStatusVal, ok := parseAttendanceStatus(newStatus)
	if !ok {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordLogConditionFail)
	}

	var oldStatusVal *ent.AttendanceStatus
	if oldStatus != nil && *oldStatus != "" {
		parsed, ok := parseAttendanceStatus(*oldStatus)
		if !ok {
			return nil, fmt.Errorf("%w", ErrAttendanceRecordLogConditionFail)
		}
		oldStatusVal = &parsed
	}

	changedByVal, err := parseOptionalUUID(changedBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordLogConditionFail)
	}

	changedAtVal := time.Now()
	if changedAt != nil && *changedAt != "" {
		parsed, err := time.Parse(time.RFC3339, *changedAt)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrAttendanceRecordLogConditionFail)
		}
		changedAtVal = parsed
	}

	item, err := s.db.CreateAttendanceRecordLog(ctx, &ent.AttendanceRecordLog{
		RecordID:  recordID,
		OldStatus: oldStatusVal,
		NewStatus: newStatusVal,
		ChangedBy: changedByVal,
		ChangedAt: changedAtVal,
		Reason:    reason,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

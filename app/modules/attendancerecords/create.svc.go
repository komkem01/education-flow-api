package attendancerecords

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, sessionID uuid.UUID, enrollmentID uuid.UUID, status string, source string, markedAt *string, remark *string, markedBy *string) (*ent.AttendanceRecord, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "attendancerecords.service.create")
	defer span.End()

	statusVal, ok := parseAttendanceStatus(status)
	if !ok {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
	}
	sourceVal, ok := parseAttendanceSource(source)
	if !ok {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
	}

	markedByVal, err := parseOptionalUUID(markedBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
	}

	markedAtVal := time.Now()
	if markedAt != nil && *markedAt != "" {
		parsed, err := time.Parse(time.RFC3339, *markedAt)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
		}
		markedAtVal = parsed
	}

	item, err := s.db.CreateAttendanceRecord(ctx, &ent.AttendanceRecord{
		SessionID:    sessionID,
		EnrollmentID: enrollmentID,
		Status:       statusVal,
		Source:       sourceVal,
		MarkedAt:     markedAtVal,
		Remark:       remark,
		MarkedBy:     markedByVal,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

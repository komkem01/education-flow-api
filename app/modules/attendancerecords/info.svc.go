package attendancerecords

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.AttendanceRecord, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "attendancerecords.service.info")
	defer span.End()

	item, err := s.db.GetAttendanceRecordByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

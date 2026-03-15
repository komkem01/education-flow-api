package attendancesessions

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.AttendanceSession, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "attendancesessions.service.info")
	defer span.End()

	item, err := s.db.GetAttendanceSessionByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

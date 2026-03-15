package classrooms

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.Classroom, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "classrooms.service.get_by_id")
	defer span.End()

	item, err := s.db.GetClassroomByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

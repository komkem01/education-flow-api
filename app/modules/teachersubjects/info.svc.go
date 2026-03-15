package teachersubjects

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.TeacherSubject, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teachersubjects.service.get_by_id")
	defer span.End()

	item, err := s.db.GetTeacherSubjectByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

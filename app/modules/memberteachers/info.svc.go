package memberteachers

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.MemberTeacher, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberteachers.service.get_by_id")
	defer span.End()

	teacher, err := s.db.GetMemberTeacherByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return teacher, nil
}

package memberteachers

import (
	"context"

	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberteachers.service.delete")
	defer span.End()

	if err := s.db.SoftDeleteMemberTeacherByID(ctx, id); err != nil {
		return normalizeServiceError(err)
	}
	return nil
}

package schooldepartments

import (
	"context"

	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "schooldepartments.service.delete")
	defer span.End()

	if err := s.db.SoftDeleteSchoolDepartmentByID(ctx, id); err != nil {
		return normalizeServiceError(err)
	}
	return nil
}

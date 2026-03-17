package schooldepartments

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Info(ctx context.Context, id uuid.UUID) (*ent.SchoolDepartment, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "schooldepartments.service.info")
	defer span.End()

	item, err := s.db.GetSchoolDepartmentByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

package departments

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.Department, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "departments.service.get_by_id")
	defer span.End()

	department, err := s.db.GetDepartmentByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return department, nil
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*ent.Department, error) {
	return s.GetByID(ctx, id)
}

package departments

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
)

func (s *Service) Create(ctx context.Context, name string, isActive bool) (*ent.Department, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "departments.service.create")
	defer span.End()

	department, err := s.db.CreateDepartment(ctx, name, isActive)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return department, nil
}

func (s *Service) CreateDepartmentService(ctx context.Context, name string, isActive bool) (*ent.Department, error) {
	return s.Create(ctx, name, isActive)
}

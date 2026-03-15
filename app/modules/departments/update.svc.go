package departments

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, name *string, isActive *bool) (*ent.Department, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "departments.service.update")
	defer span.End()

	if name == nil && isActive == nil {
		return nil, fmt.Errorf("%w", ErrDepartmentConditionFail)
	}

	department, err := s.db.UpdateDepartmentByID(ctx, id, name, isActive)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return department, nil
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, name *string, isActive *bool) (*ent.Department, error) {
	return s.Update(ctx, id, name, isActive)
}

package schooldepartments

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, data *ent.SchoolDepartmentUpdate) (*ent.SchoolDepartment, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "schooldepartments.service.update")
	defer span.End()

	if data == nil {
		return nil, fmt.Errorf("%w", ErrSchoolDepartmentConditionFail)
	}

	if data.SchoolID == nil && data.DepartmentID == nil && data.Code == nil && data.CustomName == nil && data.IsActive == nil {
		return nil, fmt.Errorf("%w", ErrSchoolDepartmentConditionFail)
	}

	item, err := s.db.UpdateSchoolDepartmentByID(ctx, id, data)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

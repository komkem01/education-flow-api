package schooldepartments

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, schoolID *uuid.UUID, departmentID *uuid.UUID, isActive *bool) ([]*ent.SchoolDepartment, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "schooldepartments.service.list")
	defer span.End()

	items, page, err := s.db.ListSchoolDepartments(ctx, req, schoolID, departmentID, isActive)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrSchoolDepartmentConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

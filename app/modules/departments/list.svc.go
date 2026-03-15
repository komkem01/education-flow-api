package departments

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, isActive *bool) ([]*ent.Department, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "departments.service.list")
	defer span.End()

	departments, page, err := s.db.ListDepartments(ctx, req, isActive)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrDepartmentConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return departments, page, nil
}

func (s *Service) ListService(ctx context.Context, req *base.RequestPaginate, isActive *bool) ([]*ent.Department, *base.ResponsePaginate, error) {
	return s.List(ctx, req, isActive)
}

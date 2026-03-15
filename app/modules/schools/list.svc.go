package schools

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate) ([]*ent.School, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "schools.service.list")
	defer span.End()

	schools, page, err := s.db.ListSchools(ctx, req)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrSchoolConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return schools, page, nil
}

func (s *Service) ListService(ctx context.Context, req *base.RequestPaginate) ([]*ent.School, *base.ResponsePaginate, error) {
	return s.List(ctx, req)
}

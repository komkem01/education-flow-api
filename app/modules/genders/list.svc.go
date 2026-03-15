package genders

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, isActive *bool) ([]*ent.Gender, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "genders.service.list")
	defer span.End()

	genders, page, err := s.db.ListGenders(ctx, req, isActive)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrGenderConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return genders, page, nil
}

func (s *Service) ListService(ctx context.Context, req *base.RequestPaginate, isActive *bool) ([]*ent.Gender, *base.ResponsePaginate, error) {
	return s.List(ctx, req, isActive)
}

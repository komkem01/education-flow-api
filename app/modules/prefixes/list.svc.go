package prefixes

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, isActive *bool, genderID *uuid.UUID) ([]*ent.Prefix, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "prefixes.service.list")
	defer span.End()

	prefixes, page, err := s.db.ListPrefixes(ctx, req, isActive, genderID)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrPrefixConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return prefixes, page, nil
}

func (s *Service) ListService(ctx context.Context, req *base.RequestPaginate, isActive *bool, genderID *uuid.UUID) ([]*ent.Prefix, *base.ResponsePaginate, error) {
	return s.List(ctx, req, isActive, genderID)
}

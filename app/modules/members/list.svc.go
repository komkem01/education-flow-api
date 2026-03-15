package members

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, role *ent.MemberRole) ([]*ent.Member, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "members.service.list")
	defer span.End()

	items, page, err := s.db.ListMembers(ctx, req, isActive, schoolID, role)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrMemberConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

func (s *Service) ListService(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, role *ent.MemberRole) ([]*ent.Member, *base.ResponsePaginate, error) {
	return s.List(ctx, req, isActive, schoolID, role)
}

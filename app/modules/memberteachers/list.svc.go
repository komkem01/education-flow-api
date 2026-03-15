package memberteachers

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, isActive *bool, memberID *uuid.UUID, departmentID *uuid.UUID) ([]*ent.MemberTeacher, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberteachers.service.list")
	defer span.End()

	teachers, page, err := s.db.ListMemberTeachers(ctx, req, isActive, memberID, departmentID)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrMemberTeacherConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return teachers, page, nil
}

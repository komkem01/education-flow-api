package memberstudents

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, advisorTeacherID *uuid.UUID) ([]*ent.MemberStudent, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberstudents.service.list")
	defer span.End()

	items, page, err := s.db.ListMemberStudents(ctx, req, isActive, schoolID, advisorTeacherID)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrMemberStudentConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

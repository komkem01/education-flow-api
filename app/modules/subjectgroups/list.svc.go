package subjectgroups

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, headTeacherID *uuid.UUID) ([]*ent.SubjectGroup, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "subjectgroups.service.list")
	defer span.End()

	items, page, err := s.db.ListSubjectGroups(ctx, req, isActive, schoolID, headTeacherID)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrSubjectGroupConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

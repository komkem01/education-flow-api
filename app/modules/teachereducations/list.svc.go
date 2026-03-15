package teachereducations

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, degree *string) ([]*ent.TeacherEducation, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teachereducations.service.list")
	defer span.End()

	var parsedDegree *ent.TeacherDegree
	if degree != nil {
		v, ok := parseTeacherDegree(*degree)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrTeacherEducationConditionFail)
		}
		parsedDegree = &v
	}

	items, page, err := s.db.ListTeacherEducations(ctx, req, teacherID, parsedDegree)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrTeacherEducationConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

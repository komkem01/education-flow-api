package teachersubjects

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, subjectID *uuid.UUID, role *string, isActive *bool) ([]*ent.TeacherSubject, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teachersubjects.service.list")
	defer span.End()

	var parsedRole *ent.TeacherSubjectRole
	if role != nil {
		v, ok := parseTeacherSubjectRole(*role)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrTeacherSubjectConditionFail)
		}
		parsedRole = &v
	}

	items, page, err := s.db.ListTeacherSubjects(ctx, req, teacherID, subjectID, parsedRole, isActive)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrTeacherSubjectConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

package teachersubjects

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, teacherID uuid.UUID, subjectID uuid.UUID, role string, isActive bool) (*ent.TeacherSubject, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teachersubjects.service.create")
	defer span.End()

	parsedRole, ok := parseTeacherSubjectRole(role)
	if !ok {
		return nil, fmt.Errorf("%w", ErrTeacherSubjectConditionFail)
	}

	item, err := s.db.CreateTeacherSubject(ctx, &ent.TeacherSubject{
		TeacherID: teacherID,
		SubjectID: subjectID,
		Role:      parsedRole,
		IsActive:  isActive,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

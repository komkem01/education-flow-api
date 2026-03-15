package teachereducations

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, teacherID uuid.UUID, degree string, major string, university string, graduationYear string) (*ent.TeacherEducation, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teachereducations.service.create")
	defer span.End()

	parsedDegree, ok := parseTeacherDegree(degree)
	if !ok {
		return nil, fmt.Errorf("%w", ErrTeacherEducationConditionFail)
	}

	item, err := s.db.CreateTeacherEducation(ctx, &ent.TeacherEducation{
		TeacherID:      teacherID,
		Degree:         parsedDegree,
		Major:          major,
		University:     university,
		GraduationYear: graduationYear,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

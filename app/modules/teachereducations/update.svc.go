package teachereducations

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.TeacherEducation, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teachereducations.service.update")
	defer span.End()

	if req.TeacherID == nil && req.Degree == nil && req.Major == nil && req.University == nil && req.GraduationYear == nil {
		return nil, fmt.Errorf("%w", ErrTeacherEducationConditionFail)
	}

	payload := &ent.TeacherEducationUpdate{
		Major:          req.Major,
		University:     req.University,
		GraduationYear: req.GraduationYear,
	}

	if req.TeacherID != nil {
		parsed, err := uuid.Parse(*req.TeacherID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrTeacherEducationConditionFail)
		}
		payload.TeacherID = &parsed
	}

	if req.Degree != nil {
		parsed, ok := parseTeacherDegree(*req.Degree)
		if !ok {
			return nil, fmt.Errorf("%w", ErrTeacherEducationConditionFail)
		}
		payload.Degree = &parsed
	}

	item, err := s.db.UpdateTeacherEducationByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

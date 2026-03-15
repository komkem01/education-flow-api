package teachersubjects

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.TeacherSubject, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teachersubjects.service.update")
	defer span.End()

	if req.TeacherID == nil && req.SubjectID == nil && req.Role == nil && req.IsActive == nil {
		return nil, fmt.Errorf("%w", ErrTeacherSubjectConditionFail)
	}

	payload := &ent.TeacherSubjectUpdate{IsActive: req.IsActive}

	if req.TeacherID != nil {
		parsed, err := uuid.Parse(*req.TeacherID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrTeacherSubjectConditionFail)
		}
		payload.TeacherID = &parsed
	}
	if req.SubjectID != nil {
		parsed, err := uuid.Parse(*req.SubjectID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrTeacherSubjectConditionFail)
		}
		payload.SubjectID = &parsed
	}
	if req.Role != nil {
		parsed, ok := parseTeacherSubjectRole(*req.Role)
		if !ok {
			return nil, fmt.Errorf("%w", ErrTeacherSubjectConditionFail)
		}
		payload.Role = &parsed
	}

	item, err := s.db.UpdateTeacherSubjectByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

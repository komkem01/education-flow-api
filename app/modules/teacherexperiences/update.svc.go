package teacherexperiences

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.TeacherExperience, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacherexperiences.service.update")
	defer span.End()

	if req.TeacherID == nil && req.SchoolName == nil && req.Position == nil && req.DepartmentName == nil && req.StartDate == nil && req.EndDate == nil && req.IsCurrent == nil && req.Responsibilities == nil && req.Achievements == nil && req.SortOrder == nil && req.IsActive == nil {
		return nil, fmt.Errorf("%w", ErrTeacherExperienceConditionFail)
	}

	payload := &ent.TeacherExperienceUpdate{
		SchoolName:       req.SchoolName,
		Position:         req.Position,
		DepartmentName:   req.DepartmentName,
		IsCurrent:        req.IsCurrent,
		Responsibilities: req.Responsibilities,
		Achievements:     req.Achievements,
		SortOrder:        req.SortOrder,
		IsActive:         req.IsActive,
	}

	if req.TeacherID != nil {
		parsed, err := uuid.Parse(*req.TeacherID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrTeacherExperienceConditionFail)
		}
		payload.TeacherID = &parsed
	}
	if req.StartDate != nil {
		t, err := parseDate(*req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrTeacherExperienceConditionFail, err)
		}
		payload.StartDate = &t
	}
	if req.EndDate != nil {
		t, err := parseOptionalDate(req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrTeacherExperienceConditionFail, err)
		}
		payload.EndDate = t
	}

	if payload.StartDate != nil && payload.EndDate != nil && payload.EndDate.Before(*payload.StartDate) {
		return nil, fmt.Errorf("%w", ErrTeacherExperienceConditionFail)
	}

	item, err := s.db.UpdateTeacherExperienceByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

package teacherexperiences

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func parseDate(v string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid-date-format: %w", err)
	}
	return t, nil
}

func parseOptionalDate(v *string) (*time.Time, error) {
	if v == nil || *v == "" {
		return nil, nil
	}
	t, err := parseDate(*v)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Service) Create(ctx context.Context, teacherID uuid.UUID, schoolName string, position string, departmentName *string, startDate string, endDate *string, isCurrent bool, responsibilities *string, achievements *string, sortOrder int, isActive bool) (*ent.TeacherExperience, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacherexperiences.service.create")
	defer span.End()

	startAt, err := parseDate(startDate)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTeacherExperienceConditionFail, err)
	}
	endAt, err := parseOptionalDate(endDate)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTeacherExperienceConditionFail, err)
	}
	if endAt != nil && endAt.Before(startAt) {
		return nil, fmt.Errorf("%w", ErrTeacherExperienceConditionFail)
	}

	item, err := s.db.CreateTeacherExperience(ctx, &ent.TeacherExperience{
		TeacherID:        teacherID,
		SchoolName:       schoolName,
		Position:         position,
		DepartmentName:   departmentName,
		StartDate:        startAt,
		EndDate:          endAt,
		IsCurrent:        isCurrent,
		Responsibilities: responsibilities,
		Achievements:     achievements,
		SortOrder:        sortOrder,
		IsActive:         isActive,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

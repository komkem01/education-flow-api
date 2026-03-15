package memberteachers

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

func (s *Service) Create(ctx context.Context, memberID uuid.UUID, genderID uuid.UUID, prefixID uuid.UUID, code string, citizenID string, firstNameTH string, lastNameTH string, firstNameEN string, lastNameEN string, phone string, position string, academicStanding string, departmentID uuid.UUID, startDate string, endDate *string, isActive bool) (*ent.MemberTeacher, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberteachers.service.create")
	defer span.End()

	startAt, err := parseDate(startDate)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMemberTeacherConditionFail, err)
	}

	endAt, err := parseOptionalDate(endDate)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMemberTeacherConditionFail, err)
	}

	teacher, err := s.db.CreateMemberTeacher(ctx, &ent.MemberTeacher{
		MemberID:         memberID,
		GenderID:         genderID,
		PrefixID:         prefixID,
		Code:             code,
		CitizenID:        citizenID,
		FirstNameTH:      firstNameTH,
		LastNameTH:       lastNameTH,
		FirstNameEN:      firstNameEN,
		LastNameEN:       lastNameEN,
		Phone:            phone,
		Position:         position,
		AcademicStanding: academicStanding,
		DepartmentID:     departmentID,
		StartDate:        startAt,
		EndDate:          endAt,
		IsActive:         isActive,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return teacher, nil
}

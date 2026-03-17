package memberteachers

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/hashing"

	"github.com/google/uuid"
)

var teacherPhoneAllowed = regexp.MustCompile(`^[0-9+\-\s]{9,20}$`)

func normalizeTeacherRequired(v string) string {
	return strings.TrimSpace(v)
}

func validateTeacherEmailPassword(email string, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return "", fmt.Errorf("%w", ErrTeacherInvalidEmail)
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return "", fmt.Errorf("%w", ErrTeacherInvalidEmail)
	}
	if len(strings.TrimSpace(password)) < 8 {
		return "", fmt.Errorf("%w", ErrTeacherInvalidPassword)
	}
	return email, nil
}

func validateTeacherCitizenID(v string) (string, error) {
	v = strings.TrimSpace(v)
	if len(v) != 13 {
		return "", fmt.Errorf("%w", ErrTeacherInvalidCitizenID)
	}
	for _, ch := range v {
		if ch < '0' || ch > '9' {
			return "", fmt.Errorf("%w", ErrTeacherInvalidCitizenID)
		}
	}
	return v, nil
}

func validateTeacherPhone(v string) (string, error) {
	v = strings.TrimSpace(v)
	if !teacherPhoneAllowed.MatchString(v) {
		return "", fmt.Errorf("%w", ErrTeacherInvalidPhone)
	}
	return v, nil
}

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

func (s *Service) Create(ctx context.Context, actorRole ent.MemberRole, schoolID uuid.UUID, email string, password string, genderID uuid.UUID, prefixID uuid.UUID, citizenID string, firstNameTH string, lastNameTH string, firstNameEN string, lastNameEN string, phone string, position string, academicStanding string, departmentID uuid.UUID, startDate string, endDate *string) (*ent.TeacherRegistrationResult, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberteachers.service.register")
	defer span.End()

	var err error
	email, err = validateTeacherEmailPassword(email, password)
	if err != nil {
		return nil, err
	}
	firstNameTH = normalizeTeacherRequired(firstNameTH)
	lastNameTH = normalizeTeacherRequired(lastNameTH)
	firstNameEN = normalizeTeacherRequired(firstNameEN)
	lastNameEN = normalizeTeacherRequired(lastNameEN)
	position = normalizeTeacherRequired(position)
	academicStanding = normalizeTeacherRequired(academicStanding)
	if firstNameTH == "" || lastNameTH == "" || firstNameEN == "" || lastNameEN == "" || position == "" || academicStanding == "" {
		return nil, fmt.Errorf("%w", ErrMemberTeacherConditionFail)
	}
	citizenID, err = validateTeacherCitizenID(citizenID)
	if err != nil {
		return nil, err
	}
	phone, err = validateTeacherPhone(phone)
	if err != nil {
		return nil, err
	}

	if actorRole != ent.MemberRoleAdmin && actorRole != ent.MemberRoleSuperadmin {
		return nil, fmt.Errorf("%w", ErrMemberTeacherUnauthorized)
	}

	if strings.TrimSpace(password) == "" {
		return nil, fmt.Errorf("%w", ErrMemberTeacherConditionFail)
	}

	startAt, err := parseDate(startDate)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMemberTeacherConditionFail, err)
	}

	endAt, err := parseOptionalDate(endDate)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMemberTeacherConditionFail, err)
	}
	if endAt != nil && endAt.Before(startAt) {
		return nil, fmt.Errorf("%w", ErrTeacherInvalidDateRange)
	}

	hashed, err := hashing.HashPassword(password)
	if err != nil {
		return nil, err
	}

	code, err := utils.GenerateNumericCode("TCH", 6)
	if err != nil {
		return nil, err
	}

	result, err := s.db.RegisterTeacher(ctx, &ent.TeacherRegistrationInput{
		MemberEmail:             email,
		MemberPasswordHash:      string(hashed),
		MemberSchoolID:          schoolID,
		MemberRole:              ent.MemberRoleTeacher,
		MemberIsActive:          true,
		MemberLastLogin:         nil,
		TeacherGenderID:         genderID,
		TeacherPrefixID:         prefixID,
		TeacherCode:             code,
		TeacherCitizenID:        citizenID,
		TeacherFirstNameTH:      firstNameTH,
		TeacherLastNameTH:       lastNameTH,
		TeacherFirstNameEN:      firstNameEN,
		TeacherLastNameEN:       lastNameEN,
		TeacherPhone:            phone,
		TeacherPosition:         position,
		TeacherAcademicStanding: academicStanding,
		TeacherDepartmentID:     departmentID,
		TeacherStartDate:        startAt,
		TeacherEndDate:          endAt,
		TeacherIsActive:         true,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return result, nil
}

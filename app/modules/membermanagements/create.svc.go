package membermanagements

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

var managementPhoneAllowed = regexp.MustCompile(`^[0-9+\-\s]{9,20}$`)

func normalizeManagementRequired(v string) string {
	return strings.TrimSpace(v)
}

func isEmployeeCodeDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "employee_code")
}

func normalizeManagementOptional(v *string) *string {
	if v == nil {
		return nil
	}
	t := strings.TrimSpace(*v)
	if t == "" {
		return nil
	}
	return &t
}

func normalizeManagementName(v *string) *string {
	if v == nil {
		return nil
	}
	t := strings.TrimSpace(*v)
	if t == "" {
		return nil
	}
	return &t
}

func validateManagementPhone(v *string) (*string, error) {
	v = normalizeManagementOptional(v)
	if v == nil {
		return nil, nil
	}
	if !managementPhoneAllowed.MatchString(*v) {
		return nil, fmt.Errorf("%w", ErrManagementInvalidPhone)
	}
	return v, nil
}

func validateManagementEmailPassword(email string, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return "", fmt.Errorf("%w", ErrManagementInvalidEmail)
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return "", fmt.Errorf("%w", ErrManagementInvalidEmail)
	}
	if len(strings.TrimSpace(password)) < 8 {
		return "", fmt.Errorf("%w", ErrManagementInvalidPassword)
	}
	return email, nil
}

func parseDate(v string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid-date-format: %w", err)
	}
	return t, nil
}

func (s *Service) resolveSchoolDepartment(ctx context.Context, schoolID uuid.UUID, schoolDepartmentID *uuid.UUID, departmentID *uuid.UUID) (*ent.SchoolDepartment, error) {
	if schoolDepartmentID != nil {
		item, err := s.schoolDepartmentDB.GetSchoolDepartmentByID(ctx, *schoolDepartmentID)
		if err != nil {
			return nil, normalizeServiceError(err)
		}
		if item.SchoolID != schoolID {
			return nil, fmt.Errorf("%w", ErrMemberManagementUnauthorized)
		}
		return item, nil
	}

	if departmentID == nil {
		return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
	}

	item, err := s.schoolDepartmentDB.GetSchoolDepartmentBySchoolAndDepartment(ctx, schoolID, *departmentID)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

func (s *Service) Create(ctx context.Context, actorID uuid.UUID, actorRole ent.MemberRole, schoolID uuid.UUID, email string, password string, genderID *uuid.UUID, prefixID *uuid.UUID, firstName *string, lastName *string, phone *string, position string, startWorkDate string, schoolDepartmentID *uuid.UUID, departmentID *uuid.UUID, requestReason *string) (*ent.ManagementRegistrationResult, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "membermanagements.service.register")
	defer span.End()

	var err error
	email, err = validateManagementEmailPassword(email, password)
	if err != nil {
		return nil, err
	}
	position = normalizeManagementRequired(position)
	if position == "" {
		return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
	}
	firstName = normalizeManagementName(firstName)
	lastName = normalizeManagementName(lastName)
	phone, err = validateManagementPhone(phone)
	if err != nil {
		return nil, err
	}
	requestReason = normalizeManagementOptional(requestReason)

	if actorRole != ent.MemberRoleAdmin && actorRole != ent.MemberRoleSuperadmin {
		return nil, fmt.Errorf("%w", ErrMemberManagementUnauthorized)
	}

	if strings.TrimSpace(password) == "" {
		return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
	}

	schoolDepartment, err := s.resolveSchoolDepartment(ctx, schoolID, schoolDepartmentID, departmentID)
	if err != nil {
		return nil, err
	}

	startAt, err := parseDate(strings.TrimSpace(startWorkDate))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMemberManagementConditionFail, err)
	}

	hashed, err := hashing.HashPassword(password)
	if err != nil {
		return nil, err
	}

	const maxEmployeeCodeRetry = 10
	for i := 0; i < maxEmployeeCodeRetry; i++ {
		employeeCode, genErr := utils.GenerateNumericCode("EMP", 6)
		if genErr != nil {
			return nil, genErr
		}

		item, registerErr := s.db.RegisterManagement(ctx, &ent.ManagementRegistrationInput{
			MemberEmail:                  email,
			MemberPasswordHash:           string(hashed),
			MemberSchoolID:               schoolID,
			MemberRole:                   ent.MemberRoleStaff,
			MemberIsActive:               false,
			MemberLastLogin:              nil,
			ManagementEmployeeCode:       employeeCode,
			ManagementGenderID:           genderID,
			ManagementPrefixID:           prefixID,
			ManagementFirstName:          firstName,
			ManagementLastName:           lastName,
			ManagementPhone:              phone,
			ManagementPosition:           position,
			ManagementStartWorkDate:      startAt,
			ManagementSchoolDepartmentID: schoolDepartment.ID,
			ManagementDepartmentID:       schoolDepartment.DepartmentID,
			ManagementIsActive:           false,
			RequestedBy:                  actorID,
			RequestedByRole:              ent.ApprovalActorRoleAdmin,
			RequestReason:                requestReason,
		})
		if registerErr == nil {
			return item, nil
		}
		if isDuplicateKeyError(registerErr) && isEmployeeCodeDuplicateError(registerErr) {
			continue
		}
		return nil, normalizeServiceError(registerErr)
	}

	return nil, fmt.Errorf("%w", ErrMemberManagementDuplicate)
}

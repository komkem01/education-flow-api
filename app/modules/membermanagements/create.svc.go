package membermanagements

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/hashing"

	"github.com/google/uuid"
)

func normalizeManagementRequired(v string) string {
	return strings.TrimSpace(v)
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

func (s *Service) Create(ctx context.Context, actorID uuid.UUID, actorRole ent.MemberRole, schoolID uuid.UUID, email string, password string, employeeCode string, position string, startWorkDate string, departmentID uuid.UUID, requestReason *string) (*ent.ManagementRegistrationResult, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "membermanagements.service.register")
	defer span.End()

	var err error
	email, err = validateManagementEmailPassword(email, password)
	if err != nil {
		return nil, err
	}
	employeeCode = normalizeManagementRequired(employeeCode)
	position = normalizeManagementRequired(position)
	if employeeCode == "" || position == "" {
		return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
	}
	requestReason = normalizeManagementOptional(requestReason)
	if requestReason == nil {
		return nil, fmt.Errorf("%w", ErrManagementInvalidReason)
	}

	if actorRole != ent.MemberRoleAdmin && actorRole != ent.MemberRoleSuperadmin {
		return nil, fmt.Errorf("%w", ErrMemberManagementUnauthorized)
	}

	if strings.TrimSpace(password) == "" {
		return nil, fmt.Errorf("%w", ErrMemberManagementConditionFail)
	}

	startAt, err := parseDate(strings.TrimSpace(startWorkDate))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMemberManagementConditionFail, err)
	}

	hashed, err := hashing.HashPassword(password)
	if err != nil {
		return nil, err
	}

	item, err := s.db.RegisterManagement(ctx, &ent.ManagementRegistrationInput{
		MemberEmail:             email,
		MemberPasswordHash:      string(hashed),
		MemberSchoolID:          schoolID,
		MemberRole:              ent.MemberRoleAdmin,
		MemberIsActive:          false,
		MemberLastLogin:         nil,
		ManagementEmployeeCode:  employeeCode,
		ManagementPosition:      position,
		ManagementStartWorkDate: startAt,
		ManagementDepartmentID:  departmentID,
		ManagementIsActive:      false,
		RequestedBy:             actorID,
		RequestedByRole:         ent.ApprovalActorRoleAdmin,
		RequestReason:           requestReason,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

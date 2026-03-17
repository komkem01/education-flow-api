package memberstudents

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

var studentPhoneAllowed = regexp.MustCompile(`^[0-9+\-\s]{9,20}$`)

func normalizeRequired(v string) string {
	return strings.TrimSpace(v)
}

func isStudentCodeDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "ux_member_students_school_student_code") || strings.Contains(errStr, "student_code")
}

func normalizeOptional(v *string) *string {
	if v == nil {
		return nil
	}
	t := strings.TrimSpace(*v)
	if t == "" {
		return nil
	}
	return &t
}

func validateEmailAndPassword(email string, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return "", fmt.Errorf("%w", ErrInvalidEmail)
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return "", fmt.Errorf("%w", ErrInvalidEmail)
	}
	if len(strings.TrimSpace(password)) < 8 {
		return "", fmt.Errorf("%w", ErrInvalidPassword)
	}
	return email, nil
}

func validateOptionalCitizenID(v *string) (*string, error) {
	v = normalizeOptional(v)
	if v == nil {
		return nil, nil
	}
	if len(*v) != 13 {
		return nil, fmt.Errorf("%w", ErrInvalidCitizenID)
	}
	for _, ch := range *v {
		if ch < '0' || ch > '9' {
			return nil, fmt.Errorf("%w", ErrInvalidCitizenID)
		}
	}
	return v, nil
}

func validateOptionalPhone(v *string) (*string, error) {
	v = normalizeOptional(v)
	if v == nil {
		return nil, nil
	}
	if !studentPhoneAllowed.MatchString(*v) {
		return nil, fmt.Errorf("%w", ErrInvalidPhone)
	}
	return v, nil
}

func (s *Service) Create(ctx context.Context, actorID uuid.UUID, actorRole ent.MemberRole, schoolID uuid.UUID, email string, password string, genderID uuid.UUID, prefixID uuid.UUID, advisorTeacherID *uuid.UUID, firstNameTH string, lastNameTH string, firstNameEN *string, lastNameEN *string, citizenID *string, phone *string, birthDate *string, nationality *string, religion *string, addressCurrent *string, addressRegistered *string, emergencyContactName *string, emergencyContactPhone *string, requestReason *string) (*ent.StudentRegistrationResult, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberstudents.service.register")
	defer span.End()

	var err error
	email, err = validateEmailAndPassword(email, password)
	if err != nil {
		return nil, err
	}

	firstNameTH = normalizeRequired(firstNameTH)
	lastNameTH = normalizeRequired(lastNameTH)
	if firstNameTH == "" || lastNameTH == "" {
		return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
	}

	firstNameEN = normalizeOptional(firstNameEN)
	lastNameEN = normalizeOptional(lastNameEN)
	if (firstNameEN == nil) != (lastNameEN == nil) {
		return nil, fmt.Errorf("%w", ErrInvalidNamePair)
	}

	citizenID, err = validateOptionalCitizenID(citizenID)
	if err != nil {
		return nil, err
	}
	phone, err = validateOptionalPhone(phone)
	if err != nil {
		return nil, err
	}
	nationality = normalizeOptional(nationality)
	religion = normalizeOptional(religion)
	addressCurrent = normalizeOptional(addressCurrent)
	addressRegistered = normalizeOptional(addressRegistered)
	emergencyContactName = normalizeOptional(emergencyContactName)
	emergencyContactPhone, err = validateOptionalPhone(emergencyContactPhone)
	if err != nil {
		return nil, err
	}
	requestReason = normalizeOptional(requestReason)

	if strings.TrimSpace(password) == "" {
		return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
	}

	requireApproval := true
	requestedByRole := ent.ApprovalActorRoleAdmin
	switch actorRole {
	case ent.MemberRoleSuperadmin, ent.MemberRoleAdmin:
		requestedByRole = ent.ApprovalActorRoleAdmin
	case ent.MemberRoleTeacher:
		requestedByRole = ent.ApprovalActorRoleTeacher
	default:
		return nil, fmt.Errorf("%w", ErrMemberStudentUnauthorized)
	}

	memberActive := !requireApproval
	studentActive := !requireApproval

	hashed, err := hashing.HashPassword(password)
	if err != nil {
		return nil, err
	}

	var birthDateVal *time.Time
	if birthDate != nil && *birthDate != "" {
		parsed, err := time.Parse("2006-01-02", *birthDate)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberStudentConditionFail)
		}
		if parsed.After(time.Now()) {
			return nil, fmt.Errorf("%w", ErrInvalidBirthDate)
		}
		birthDateVal = &parsed
	}

	const maxStudentCodeRetry = 10
	for i := 0; i < maxStudentCodeRetry; i++ {
		studentCode, genErr := utils.GenerateNumericCode("STD", 6)
		if genErr != nil {
			return nil, genErr
		}

		item, registerErr := s.db.RegisterStudent(ctx, &ent.StudentRegistrationInput{
			MemberEmail:                  email,
			MemberPasswordHash:           string(hashed),
			MemberSchoolID:               schoolID,
			MemberRole:                   ent.MemberRoleStudent,
			MemberIsActive:               memberActive,
			MemberLastLogin:              nil,
			StudentSchoolID:              schoolID,
			StudentGenderID:              genderID,
			StudentPrefixID:              prefixID,
			StudentAdvisorTeacherID:      advisorTeacherID,
			StudentCode:                  studentCode,
			StudentFirstNameTH:           firstNameTH,
			StudentLastNameTH:            lastNameTH,
			StudentFirstNameEN:           firstNameEN,
			StudentLastNameEN:            lastNameEN,
			StudentCitizenID:             citizenID,
			StudentPhone:                 phone,
			StudentIsActive:              studentActive,
			ProfileBirthDate:             birthDateVal,
			ProfileNationality:           nationality,
			ProfileReligion:              religion,
			ProfileAddressCurrent:        addressCurrent,
			ProfileAddressRegistered:     addressRegistered,
			ProfileEmergencyContactName:  emergencyContactName,
			ProfileEmergencyContactPhone: emergencyContactPhone,
			RequireApproval:              requireApproval,
			RequestedBy:                  actorID,
			RequestedByRole:              requestedByRole,
			RequestReason:                requestReason,
		})
		if registerErr == nil {
			return item, nil
		}
		if isDuplicateKeyError(registerErr) && isStudentCodeDuplicateError(registerErr) {
			continue
		}
		return nil, normalizeServiceError(registerErr)
	}

	return nil, fmt.Errorf("%w", ErrMemberStudentDuplicate)
}

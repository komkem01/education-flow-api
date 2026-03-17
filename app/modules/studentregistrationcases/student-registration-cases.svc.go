package studentregistrationcases

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
)

var (
	ErrStudentRegistrationCaseNotFound      = errors.New("student-registration-case-not-found")
	ErrStudentRegistrationCaseDuplicate     = errors.New("student-registration-case-duplicate")
	ErrStudentRegistrationCaseUnauthorized  = errors.New("student-registration-case-unauthorized")
	ErrStudentRegistrationCaseConditionFail = errors.New("student-registration-case-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrStudentRegistrationCaseNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrStudentRegistrationCaseDuplicate, err)
	}
	errText := strings.ToLower(err.Error())
	if strings.Contains(errText, "invalid-status") || strings.Contains(errText, "missing-student-id") || strings.Contains(errText, "missing-student-credentials") {
		return fmt.Errorf("%w: %v", ErrStudentRegistrationCaseConditionFail, err)
	}
	if strings.Contains(errText, "foreign key constraint") ||
		strings.Contains(errText, "violates check constraint") ||
		strings.Contains(errText, "invalid input value for enum") ||
		strings.Contains(errText, "null value in column") ||
		strings.Contains(errText, "sqlstate=23503") ||
		strings.Contains(errText, "sqlstate=23514") ||
		strings.Contains(errText, "sqlstate=22p02") {
		return fmt.Errorf("%w: %v", ErrStudentRegistrationCaseConditionFail, err)
	}
	return err
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "sqlstate=23505") || strings.Contains(errStr, "duplicate key value") || strings.Contains(errStr, "violates unique constraint")
}

func parseRegistrationType(v string) (ent.StudentRegistrationType, bool) {
	t := ent.StudentRegistrationType(strings.TrimSpace(v))
	switch t {
	case ent.StudentRegistrationTypeNewEnrollment,
		ent.StudentRegistrationTypeTransferIn,
		ent.StudentRegistrationTypeTransferOut,
		ent.StudentRegistrationTypeLeaveAbsence,
		ent.StudentRegistrationTypeWithdrawal,
		ent.StudentRegistrationTypeReEnrollment:
		return t, true
	default:
		return "", false
	}
}

func parseRegistrationStatus(v string) (ent.StudentRegistrationCaseStatus, bool) {
	s := ent.StudentRegistrationCaseStatus(strings.TrimSpace(v))
	switch s {
	case ent.StudentRegistrationCaseStatusDraft,
		ent.StudentRegistrationCaseStatusPending,
		ent.StudentRegistrationCaseStatusApproved,
		ent.StudentRegistrationCaseStatusRejected,
		ent.StudentRegistrationCaseStatusCancelled,
		ent.StudentRegistrationCaseStatusApplied:
		return s, true
	default:
		return "", false
	}
}

func parseApprovalActorRoleFromMember(role ent.MemberRole) (ent.ApprovalActorRole, bool) {
	switch role {
	case ent.MemberRoleSuperadmin, ent.MemberRoleAdmin:
		return ent.ApprovalActorRoleAdmin, true
	case ent.MemberRoleTeacher:
		return ent.ApprovalActorRoleTeacher, true
	default:
		return "", false
	}
}

func parseAddressType(v string) (ent.RegistrationAddressType, bool) {
	t := ent.RegistrationAddressType(strings.TrimSpace(v))
	switch t {
	case ent.RegistrationAddressTypeCurrent, ent.RegistrationAddressTypeRegistered, ent.RegistrationAddressTypeContact:
		return t, true
	default:
		return "", false
	}
}

func parseRegistrationDocumentType(v string) (ent.RegistrationDocumentType, bool) {
	t := ent.RegistrationDocumentType(strings.TrimSpace(v))
	switch t {
	case ent.RegistrationDocumentTypeTranscript,
		ent.RegistrationDocumentTypeTransferLetter,
		ent.RegistrationDocumentTypeHouseholdReg,
		ent.RegistrationDocumentTypeBirthCertificate,
		ent.RegistrationDocumentTypeIDCard,
		ent.RegistrationDocumentTypeMedicalCert,
		ent.RegistrationDocumentTypePhoto,
		ent.RegistrationDocumentTypeOther:
		return t, true
	default:
		return "", false
	}
}

func parseIncomeBracket(v string) (ent.RegistrationIncomeBracket, bool) {
	b := ent.RegistrationIncomeBracket(strings.TrimSpace(v))
	switch b {
	case ent.RegistrationIncomeBracketUnder5000,
		ent.RegistrationIncomeBracket5001_10000,
		ent.RegistrationIncomeBracket10001_20000,
		ent.RegistrationIncomeBracket20001_40000,
		ent.RegistrationIncomeBracket40001_60000,
		ent.RegistrationIncomeBracketAbove60000:
		return b, true
	default:
		return "", false
	}
}

func normalizeStrPtr(v *string) *string {
	if v == nil {
		return nil
	}
	t := strings.TrimSpace(*v)
	if t == "" {
		return nil
	}
	return &t
}

func parseDatePtr(v *string) (*time.Time, error) {
	v = normalizeStrPtr(v)
	if v == nil {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *v)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

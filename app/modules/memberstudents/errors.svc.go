package memberstudents

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMemberStudentNotFound      = errors.New("member-student-not-found")
	ErrMemberStudentDuplicate     = errors.New("member-student-duplicate")
	ErrMemberStudentUnauthorized  = errors.New("member-student-unauthorized")
	ErrMemberStudentConditionFail = errors.New("member-student-condition-fail")
	ErrInvalidEmail               = errors.New("invalid-email")
	ErrInvalidPassword            = errors.New("invalid-password")
	ErrInvalidCitizenID           = errors.New("invalid-citizen-id")
	ErrInvalidPhone               = errors.New("invalid-phone")
	ErrInvalidBirthDate           = errors.New("invalid-birth-date")
	ErrInvalidNamePair            = errors.New("invalid-name-pair")
	ErrInvalidApprovalReason      = errors.New("invalid-approval-reason")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrMemberStudentNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrMemberStudentDuplicate, err)
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

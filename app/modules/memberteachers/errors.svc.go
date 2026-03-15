package memberteachers

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMemberTeacherNotFound      = errors.New("member-teacher-not-found")
	ErrMemberTeacherDuplicate     = errors.New("member-teacher-duplicate")
	ErrMemberTeacherUnauthorized  = errors.New("member-teacher-unauthorized")
	ErrMemberTeacherConditionFail = errors.New("member-teacher-condition-fail")
	ErrTeacherInvalidEmail        = errors.New("invalid-email")
	ErrTeacherInvalidPassword     = errors.New("invalid-password")
	ErrTeacherInvalidCitizenID    = errors.New("invalid-citizen-id")
	ErrTeacherInvalidPhone        = errors.New("invalid-phone")
	ErrTeacherInvalidDateRange    = errors.New("invalid-date-range")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrMemberTeacherNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrMemberTeacherDuplicate, err)
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

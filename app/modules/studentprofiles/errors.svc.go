package studentprofiles

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrStudentProfileNotFound      = errors.New("student-profile-not-found")
	ErrStudentProfileDuplicate     = errors.New("student-profile-duplicate")
	ErrStudentProfileUnauthorized  = errors.New("student-profile-unauthorized")
	ErrStudentProfileConditionFail = errors.New("student-profile-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrStudentProfileNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrStudentProfileDuplicate, err)
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

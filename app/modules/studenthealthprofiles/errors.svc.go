package studenthealthprofiles

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrStudentHealthProfileNotFound      = errors.New("student-health-profile-not-found")
	ErrStudentHealthProfileDuplicate     = errors.New("student-health-profile-duplicate")
	ErrStudentHealthProfileUnauthorized  = errors.New("student-health-profile-unauthorized")
	ErrStudentHealthProfileConditionFail = errors.New("student-health-profile-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrStudentHealthProfileNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrStudentHealthProfileDuplicate, err)
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

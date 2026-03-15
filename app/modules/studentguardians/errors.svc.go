package studentguardians

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrStudentGuardianNotFound      = errors.New("student-guardian-not-found")
	ErrStudentGuardianDuplicate     = errors.New("student-guardian-duplicate")
	ErrStudentGuardianUnauthorized  = errors.New("student-guardian-unauthorized")
	ErrStudentGuardianConditionFail = errors.New("student-guardian-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrStudentGuardianNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrStudentGuardianDuplicate, err)
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

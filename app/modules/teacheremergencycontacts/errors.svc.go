package teacheremergencycontacts

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrTeacherEmergencyContactNotFound      = errors.New("teacher-emergency-contact-not-found")
	ErrTeacherEmergencyContactDuplicate     = errors.New("teacher-emergency-contact-duplicate")
	ErrTeacherEmergencyContactUnauthorized  = errors.New("teacher-emergency-contact-unauthorized")
	ErrTeacherEmergencyContactConditionFail = errors.New("teacher-emergency-contact-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrTeacherEmergencyContactNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrTeacherEmergencyContactDuplicate, err)
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

func parseBoolFilter(v string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "true", "1", "t", "yes", "y":
		return true, true
	case "false", "0", "f", "no", "n":
		return false, true
	default:
		return false, false
	}
}

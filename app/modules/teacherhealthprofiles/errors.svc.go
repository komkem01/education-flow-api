package teacherhealthprofiles

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrTeacherHealthProfileNotFound      = errors.New("teacher-health-profile-not-found")
	ErrTeacherHealthProfileDuplicate     = errors.New("teacher-health-profile-duplicate")
	ErrTeacherHealthProfileUnauthorized  = errors.New("teacher-health-profile-unauthorized")
	ErrTeacherHealthProfileConditionFail = errors.New("teacher-health-profile-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrTeacherHealthProfileNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrTeacherHealthProfileDuplicate, err)
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

func parseBloodType(v string) (string, bool) {
	switch strings.ToUpper(v) {
	case "A", "B", "AB", "O":
		return strings.ToUpper(v), true
	default:
		return "", false
	}
}

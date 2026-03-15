package teacherlicenses

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"eduflow/app/modules/entities/ent"
)

var (
	ErrTeacherLicenseNotFound      = errors.New("teacher-license-not-found")
	ErrTeacherLicenseDuplicate     = errors.New("teacher-license-duplicate")
	ErrTeacherLicenseUnauthorized  = errors.New("teacher-license-unauthorized")
	ErrTeacherLicenseConditionFail = errors.New("teacher-license-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrTeacherLicenseNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrTeacherLicenseDuplicate, err)
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

func parseTeacherLicenseStatus(v string) (ent.TeacherLicenseStatus, bool) {
	status := ent.TeacherLicenseStatus(v)
	switch status {
	case ent.TeacherLicenseStatusActive, ent.TeacherLicenseStatusSuspended, ent.TeacherLicenseStatusExpired, ent.TeacherLicenseStatusRevoked:
		return status, true
	default:
		return "", false
	}
}

package enrollmentsubjects

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"eduflow/app/modules/entities/ent"
)

var (
	ErrEnrollmentSubjectNotFound      = errors.New("enrollment-subject-not-found")
	ErrEnrollmentSubjectDuplicate     = errors.New("enrollment-subject-duplicate")
	ErrEnrollmentSubjectUnauthorized  = errors.New("enrollment-subject-unauthorized")
	ErrEnrollmentSubjectConditionFail = errors.New("enrollment-subject-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrEnrollmentSubjectNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrEnrollmentSubjectDuplicate, err)
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

func parseStudentEnrollmentStatus(v string) (ent.StudentEnrollmentStatus, bool) {
	status := ent.StudentEnrollmentStatus(v)
	switch status {
	case ent.StudentEnrollmentStatusActive, ent.StudentEnrollmentStatusTransferred, ent.StudentEnrollmentStatusGraduated, ent.StudentEnrollmentStatusDropped:
		return status, true
	default:
		return "", false
	}
}

func parseBoolFilter(v string) (bool, bool) {
	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return false, false
	}
	return parsed, true
}

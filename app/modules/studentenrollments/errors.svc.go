package studentenrollments

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"eduflow/app/modules/entities/ent"
)

var (
	ErrStudentEnrollmentNotFound      = errors.New("student-enrollment-not-found")
	ErrStudentEnrollmentDuplicate     = errors.New("student-enrollment-duplicate")
	ErrStudentEnrollmentUnauthorized  = errors.New("student-enrollment-unauthorized")
	ErrStudentEnrollmentConditionFail = errors.New("student-enrollment-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrStudentEnrollmentNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrStudentEnrollmentDuplicate, err)
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

func parseEnrollmentType(v string) (ent.EnrollmentType, bool) {
	typ := ent.EnrollmentType(v)
	switch typ {
	case ent.EnrollmentTypeNew, ent.EnrollmentTypeTransferIn, ent.EnrollmentTypeRepeat, ent.EnrollmentTypeReturn:
		return typ, true
	default:
		return "", false
	}
}

func parseEnrollmentExitReason(v string) (ent.EnrollmentExitReason, bool) {
	reason := ent.EnrollmentExitReason(v)
	switch reason {
	case ent.EnrollmentExitReasonTransferOut, ent.EnrollmentExitReasonGraduated, ent.EnrollmentExitReasonDropped, ent.EnrollmentExitReasonLeave:
		return reason, true
	default:
		return "", false
	}
}

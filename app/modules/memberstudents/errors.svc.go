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

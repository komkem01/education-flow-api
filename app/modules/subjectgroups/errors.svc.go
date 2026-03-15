package subjectgroups

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrSubjectGroupNotFound      = errors.New("subject-group-not-found")
	ErrSubjectGroupDuplicate     = errors.New("subject-group-duplicate")
	ErrSubjectGroupUnauthorized  = errors.New("subject-group-unauthorized")
	ErrSubjectGroupConditionFail = errors.New("subject-group-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrSubjectGroupNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrSubjectGroupDuplicate, err)
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

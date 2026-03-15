package classrooms

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrClassroomNotFound      = errors.New("classroom-not-found")
	ErrClassroomDuplicate     = errors.New("classroom-duplicate")
	ErrClassroomUnauthorized  = errors.New("classroom-unauthorized")
	ErrClassroomConditionFail = errors.New("classroom-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrClassroomNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrClassroomDuplicate, err)
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

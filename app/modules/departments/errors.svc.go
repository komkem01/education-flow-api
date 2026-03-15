package departments

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrDepartmentNotFound      = errors.New("department-not-found")
	ErrDepartmentNameDuplicate = errors.New("department-name-duplicate")
	ErrDepartmentInvalidUpdate = errors.New("empty-update-payload")
	ErrDepartmentUnauthorized  = errors.New("department-unauthorized")
	ErrDepartmentConditionFail = errors.New("department-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrDepartmentNotFound, err)
	}

	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrDepartmentNameDuplicate, err)
	}

	return err
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "sqlstate=23505") ||
		strings.Contains(errStr, "duplicate key value") ||
		strings.Contains(errStr, "violates unique constraint")
}

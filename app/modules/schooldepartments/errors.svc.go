package schooldepartments

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrSchoolDepartmentNotFound      = errors.New("school-department-not-found")
	ErrSchoolDepartmentDuplicate     = errors.New("school-department-duplicate")
	ErrSchoolDepartmentConditionFail = errors.New("school-department-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrSchoolDepartmentNotFound, err)
	}

	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrSchoolDepartmentDuplicate, err)
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

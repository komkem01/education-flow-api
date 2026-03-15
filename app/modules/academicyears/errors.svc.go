package academicyears

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrAcademicYearNotFound      = errors.New("academic-year-not-found")
	ErrAcademicYearDuplicate     = errors.New("academic-year-duplicate")
	ErrAcademicYearInvalidUpdate = errors.New("empty-update-payload")
	ErrAcademicYearUnauthorized  = errors.New("academic-year-unauthorized")
	ErrAcademicYearConditionFail = errors.New("academic-year-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrAcademicYearNotFound, err)
	}

	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrAcademicYearDuplicate, err)
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

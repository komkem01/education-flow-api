package schools

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrSchoolNotFound      = errors.New("school-not-found")
	ErrSchoolNameDuplicate = errors.New("school-name-duplicate")
	ErrSchoolInvalidUpdate = errors.New("empty-update-payload")
	ErrSchoolUnauthorized  = errors.New("school-unauthorized")
	ErrSchoolConditionFail = errors.New("school-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrSchoolNotFound, err)
	}

	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrSchoolNameDuplicate, err)
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

package genders

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrGenderNotFound      = errors.New("gender-not-found")
	ErrGenderNameDuplicate = errors.New("gender-name-duplicate")
	ErrGenderInvalidUpdate = errors.New("empty-update-payload")
	ErrGenderUnauthorized  = errors.New("gender-unauthorized")
	ErrGenderConditionFail = errors.New("gender-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrGenderNotFound, err)
	}

	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrGenderNameDuplicate, err)
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

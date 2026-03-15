package prefixes

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrPrefixNotFound      = errors.New("prefix-not-found")
	ErrPrefixNameDuplicate = errors.New("prefix-name-duplicate")
	ErrPrefixInvalidUpdate = errors.New("empty-update-payload")
	ErrPrefixUnauthorized  = errors.New("prefix-unauthorized")
	ErrPrefixConditionFail = errors.New("prefix-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrPrefixNotFound, err)
	}

	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrPrefixNameDuplicate, err)
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

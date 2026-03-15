package auth

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrAuthUnauthorized  = errors.New("auth-unauthorized")
	ErrAuthConditionFail = errors.New("auth-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrAuthUnauthorized, err)
	}
	return err
}

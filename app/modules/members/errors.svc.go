package members

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMemberNotFound      = errors.New("member-not-found")
	ErrMemberDuplicate     = errors.New("member-duplicate")
	ErrMemberUnauthorized  = errors.New("member-unauthorized")
	ErrMemberConditionFail = errors.New("member-condition-fail")
	ErrMemberAdminLimit    = errors.New("member-admin-limit")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrMemberNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrMemberDuplicate, err)
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

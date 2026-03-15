package membermanagements

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMemberManagementNotFound      = errors.New("member-management-not-found")
	ErrMemberManagementDuplicate     = errors.New("member-management-duplicate")
	ErrMemberManagementUnauthorized  = errors.New("member-management-unauthorized")
	ErrMemberManagementConditionFail = errors.New("member-management-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrMemberManagementNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrMemberManagementDuplicate, err)
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

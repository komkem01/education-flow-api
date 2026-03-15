package memberguardians

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMemberGuardianNotFound      = errors.New("member-guardian-not-found")
	ErrMemberGuardianDuplicate     = errors.New("member-guardian-duplicate")
	ErrMemberGuardianUnauthorized  = errors.New("member-guardian-unauthorized")
	ErrMemberGuardianConditionFail = errors.New("member-guardian-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrMemberGuardianNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrMemberGuardianDuplicate, err)
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

package auditlogs

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrAuditLogNotFound      = errors.New("audit-log-not-found")
	ErrAuditLogConditionFail = errors.New("audit-log-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrAuditLogNotFound, err)
	}
	return err
}

package storages

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"eduflow/app/modules/entities/ent"
)

var (
	ErrStorageNotFound      = errors.New("storage-not-found")
	ErrStorageDuplicate     = errors.New("storage-duplicate")
	ErrStorageConditionFail = errors.New("storage-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrStorageNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrStorageDuplicate, err)
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

func parseStorageProvider(v string) (ent.StorageProvider, bool) {
	p := ent.StorageProvider(strings.TrimSpace(v))
	switch p {
	case ent.StorageProviderS3:
		return p, true
	default:
		return "", false
	}
}

package documents

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"eduflow/app/modules/s3"
)

var (
	ErrDocumentNotFound      = errors.New("document-not-found")
	ErrDocumentDuplicate     = errors.New("document-duplicate")
	ErrDocumentUnauthorized  = errors.New("document-unauthorized")
	ErrDocumentConditionFail = errors.New("document-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrDocumentNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrDocumentDuplicate, err)
	}
	if errors.Is(err, s3.ErrS3InvalidObjectKey) {
		return err
	}
	if errors.Is(err, s3.ErrS3InvalidConfig) || errors.Is(err, s3.ErrS3PresignFailed) || errors.Is(err, s3.ErrS3UploadFailed) {
		return fmt.Errorf("%w: %v", ErrDocumentConditionFail, err)
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

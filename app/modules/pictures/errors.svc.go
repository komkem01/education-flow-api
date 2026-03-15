package pictures

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"eduflow/app/modules/s3"
)

var (
	ErrPictureNotFound      = errors.New("picture-not-found")
	ErrPictureDuplicate     = errors.New("picture-duplicate")
	ErrPictureUnauthorized  = errors.New("picture-unauthorized")
	ErrPictureConditionFail = errors.New("picture-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrPictureNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrPictureDuplicate, err)
	}
	if errors.Is(err, s3.ErrS3InvalidObjectKey) {
		return err
	}
	if errors.Is(err, s3.ErrS3InvalidConfig) || errors.Is(err, s3.ErrS3PresignFailed) {
		return fmt.Errorf("%w: %v", ErrPictureConditionFail, err)
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

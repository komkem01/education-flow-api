package s3

import "errors"

var (
	ErrS3InvalidConfig    = errors.New("s3-invalid-config")
	ErrS3InvalidObjectKey = errors.New("s3-invalid-object-key")
	ErrS3PresignFailed    = errors.New("s3-presign-failed")
	ErrS3UploadFailed     = errors.New("s3-upload-failed")
)

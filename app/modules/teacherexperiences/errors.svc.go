package teacherexperiences

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrTeacherExperienceNotFound      = errors.New("teacher-experience-not-found")
	ErrTeacherExperienceDuplicate     = errors.New("teacher-experience-duplicate")
	ErrTeacherExperienceUnauthorized  = errors.New("teacher-experience-unauthorized")
	ErrTeacherExperienceConditionFail = errors.New("teacher-experience-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrTeacherExperienceNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrTeacherExperienceDuplicate, err)
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

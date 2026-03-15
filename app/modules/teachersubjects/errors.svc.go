package teachersubjects

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"eduflow/app/modules/entities/ent"
)

var (
	ErrTeacherSubjectNotFound      = errors.New("teacher-subject-not-found")
	ErrTeacherSubjectDuplicate     = errors.New("teacher-subject-duplicate")
	ErrTeacherSubjectUnauthorized  = errors.New("teacher-subject-unauthorized")
	ErrTeacherSubjectConditionFail = errors.New("teacher-subject-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrTeacherSubjectNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrTeacherSubjectDuplicate, err)
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

func parseTeacherSubjectRole(v string) (ent.TeacherSubjectRole, bool) {
	role := ent.TeacherSubjectRole(v)
	switch role {
	case ent.TeacherSubjectRolePrimary, ent.TeacherSubjectRoleAssistant:
		return role, true
	default:
		return "", false
	}
}

package teachereducations

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"eduflow/app/modules/entities/ent"
)

var (
	ErrTeacherEducationNotFound      = errors.New("teacher-education-not-found")
	ErrTeacherEducationDuplicate     = errors.New("teacher-education-duplicate")
	ErrTeacherEducationUnauthorized  = errors.New("teacher-education-unauthorized")
	ErrTeacherEducationConditionFail = errors.New("teacher-education-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrTeacherEducationNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrTeacherEducationDuplicate, err)
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

func parseTeacherDegree(v string) (ent.TeacherDegree, bool) {
	degree := ent.TeacherDegree(v)
	switch degree {
	case ent.TeacherDegreeBachelor, ent.TeacherDegreeMaster, ent.TeacherDegreeDoctor:
		return degree, true
	default:
		return "", false
	}
}

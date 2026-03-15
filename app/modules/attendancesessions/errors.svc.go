package attendancesessions

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"eduflow/app/modules/entities/ent"
)

var (
	ErrAttendanceSessionNotFound      = errors.New("attendance-session-not-found")
	ErrAttendanceSessionDuplicate     = errors.New("attendance-session-duplicate")
	ErrAttendanceSessionUnauthorized  = errors.New("attendance-session-unauthorized")
	ErrAttendanceSessionConditionFail = errors.New("attendance-session-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrAttendanceSessionNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrAttendanceSessionDuplicate, err)
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

func parseAttendanceMode(v string) (ent.AttendanceMode, bool) {
	mode := ent.AttendanceMode(v)
	switch mode {
	case ent.AttendanceModeHomeroom, ent.AttendanceModeSubject, ent.AttendanceModeActivity:
		return mode, true
	default:
		return "", false
	}
}

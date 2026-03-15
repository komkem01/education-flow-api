package attendancerecords

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"

	"github.com/google/uuid"
)

var (
	ErrAttendanceRecordNotFound      = errors.New("attendance-record-not-found")
	ErrAttendanceRecordDuplicate     = errors.New("attendance-record-duplicate")
	ErrAttendanceRecordUnauthorized  = errors.New("attendance-record-unauthorized")
	ErrAttendanceRecordConditionFail = errors.New("attendance-record-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrAttendanceRecordNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrAttendanceRecordDuplicate, err)
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

func parseAttendanceStatus(v string) (ent.AttendanceStatus, bool) {
	status := ent.AttendanceStatus(v)
	switch status {
	case ent.AttendanceStatusPresent, ent.AttendanceStatusLate, ent.AttendanceStatusAbsent, ent.AttendanceStatusSick, ent.AttendanceStatusLeave, ent.AttendanceStatusActivity:
		return status, true
	default:
		return "", false
	}
}

func parseAttendanceSource(v string) (ent.AttendanceSource, bool) {
	src := ent.AttendanceSource(v)
	switch src {
	case ent.AttendanceSourceManual, ent.AttendanceSourceQR, ent.AttendanceSourceRFID, ent.AttendanceSourceFace, ent.AttendanceSourceAPI:
		return src, true
	default:
		return "", false
	}
}

func parseOptionalUUID(v *string) (*uuid.UUID, error) {
	if v == nil || *v == "" {
		return nil, nil
	}
	parsed, err := uuid.Parse(*v)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseOptionalDateTime(v *string) (*time.Time, error) {
	if v == nil || *v == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, *v)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

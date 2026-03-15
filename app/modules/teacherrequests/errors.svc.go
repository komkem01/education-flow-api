package teacherrequests

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"eduflow/app/modules/entities/ent"
)

var (
	ErrTeacherRequestNotFound      = errors.New("teacher-request-not-found")
	ErrTeacherRequestDuplicate     = errors.New("teacher-request-duplicate")
	ErrTeacherRequestUnauthorized  = errors.New("teacher-request-unauthorized")
	ErrTeacherRequestConditionFail = errors.New("teacher-request-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrTeacherRequestNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrTeacherRequestDuplicate, err)
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

func parseTeacherRequestType(v string) (ent.TeacherRequestType, bool) {
	t := ent.TeacherRequestType(v)
	switch t {
	case ent.TeacherRequestTypeEdit, ent.TeacherRequestTypeDelete, ent.TeacherRequestTypeOther:
		return t, true
	default:
		return "", false
	}
}

func parseTeacherRequestStatus(v string) (ent.TeacherRequestStatus, bool) {
	s := ent.TeacherRequestStatus(v)
	switch s {
	case ent.TeacherRequestStatusPending, ent.TeacherRequestStatusApproved, ent.TeacherRequestStatusRejected:
		return s, true
	default:
		return "", false
	}
}

func parseApprovalActorRole(v string) (ent.ApprovalActorRole, bool) {
	r := ent.ApprovalActorRole(v)
	switch r {
	case ent.ApprovalActorRoleTeacher, ent.ApprovalActorRoleAdmin:
		return r, true
	default:
		return "", false
	}
}

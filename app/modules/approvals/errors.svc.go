package approvals

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"eduflow/app/modules/entities/ent"
)

var (
	ErrApprovalRequestNotFound      = errors.New("approval-request-not-found")
	ErrApprovalRequestDuplicate     = errors.New("approval-request-duplicate")
	ErrApprovalRequestUnauthorized  = errors.New("approval-request-unauthorized")
	ErrApprovalRequestConditionFail = errors.New("approval-request-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrApprovalRequestNotFound, err)
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: %v", ErrApprovalRequestDuplicate, err)
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

func parseApprovalActorRole(v string) (ent.ApprovalActorRole, bool) {
	r := ent.ApprovalActorRole(v)
	switch r {
	case ent.ApprovalActorRoleTeacher, ent.ApprovalActorRoleAdmin:
		return r, true
	default:
		return "", false
	}
}

func parseApprovalStatus(v string) (ent.ApprovalRequestStatus, bool) {
	s := ent.ApprovalRequestStatus(v)
	switch s {
	case ent.ApprovalRequestStatusDraft,
		ent.ApprovalRequestStatusPending,
		ent.ApprovalRequestStatusApproved,
		ent.ApprovalRequestStatusRejected,
		ent.ApprovalRequestStatusCancelled:
		return s, true
	default:
		return "", false
	}
}

func parseApprovalActionType(v string) (ent.ApprovalActionType, bool) {
	a := ent.ApprovalActionType(v)
	switch a {
	case ent.ApprovalActionTypeSubmit,
		ent.ApprovalActionTypeApprove,
		ent.ApprovalActionTypeReject,
		ent.ApprovalActionTypeCancel,
		ent.ApprovalActionTypeComment:
		return a, true
	default:
		return "", false
	}
}

func isTerminalStatus(status ent.ApprovalRequestStatus) bool {
	switch status {
	case ent.ApprovalRequestStatusApproved, ent.ApprovalRequestStatusRejected, ent.ApprovalRequestStatusCancelled:
		return true
	default:
		return false
	}
}

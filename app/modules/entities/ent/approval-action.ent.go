package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ApprovalActionType string

const (
	ApprovalActionTypeSubmit  ApprovalActionType = "submit"
	ApprovalActionTypeApprove ApprovalActionType = "approve"
	ApprovalActionTypeReject  ApprovalActionType = "reject"
	ApprovalActionTypeCancel  ApprovalActionType = "cancel"
	ApprovalActionTypeComment ApprovalActionType = "comment"
)

type ApprovalAction struct {
	bun.BaseModel `bun:"table:approval_actions,alias:aa"`

	ID             uuid.UUID          `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	RequestID      uuid.UUID          `bun:"request_id,notnull,type:uuid"`
	Action         ApprovalActionType `bun:"action,notnull,type:approval_action_type"`
	IdempotencyKey *string            `bun:"idempotency_key"`
	ActedBy        uuid.UUID          `bun:"acted_by,notnull,type:uuid"`
	ActedByRole    ApprovalActorRole  `bun:"acted_by_role,notnull,type:approval_actor_role"`
	Comment        *string            `bun:"comment"`
	Metadata       map[string]any     `bun:"metadata,type:jsonb"`
	CreatedAt      time.Time          `bun:"created_at,notnull,default:current_timestamp"`
}

package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ApprovalActorRole string

const (
	ApprovalActorRoleTeacher ApprovalActorRole = "teacher"
	ApprovalActorRoleAdmin   ApprovalActorRole = "admin"
)

type ApprovalRequestStatus string

const (
	ApprovalRequestStatusDraft     ApprovalRequestStatus = "draft"
	ApprovalRequestStatusPending   ApprovalRequestStatus = "pending"
	ApprovalRequestStatusApproved  ApprovalRequestStatus = "approved"
	ApprovalRequestStatusRejected  ApprovalRequestStatus = "rejected"
	ApprovalRequestStatusCancelled ApprovalRequestStatus = "cancelled"
)

type ApprovalRequest struct {
	bun.BaseModel `bun:"table:approval_requests,alias:ar"`

	ID              uuid.UUID             `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	RegistrationNo  string                `bun:"registration_no,notnull,type:varchar(32)"`
	RequestType     string                `bun:"request_type,notnull"`
	SubjectType     string                `bun:"subject_type,notnull"`
	SubjectID       *uuid.UUID            `bun:"subject_id,type:uuid"`
	RequestedBy     uuid.UUID             `bun:"requested_by,notnull,type:uuid"`
	RequestedByRole ApprovalActorRole     `bun:"requested_by_role,notnull,type:approval_actor_role"`
	Payload         map[string]any        `bun:"payload,notnull,type:jsonb"`
	CurrentStatus   ApprovalRequestStatus `bun:"current_status,notnull,type:approval_request_status,default:draft"`
	SubmittedAt     *time.Time            `bun:"submitted_at"`
	ResolvedAt      *time.Time            `bun:"resolved_at"`
	CreatedAt       time.Time             `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt       time.Time             `bun:"updated_at,notnull,default:current_timestamp"`
}

type ApprovalRequestUpdate struct {
	RegistrationNo  *string
	RequestType     *string
	SubjectType     *string
	SubjectID       *uuid.UUID
	RequestedBy     *uuid.UUID
	RequestedByRole *ApprovalActorRole
	Payload         *map[string]any
	CurrentStatus   *ApprovalRequestStatus
	SubmittedAt     *time.Time
	ResolvedAt      *time.Time
}

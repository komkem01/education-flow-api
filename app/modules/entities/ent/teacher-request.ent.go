package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeacherRequestType string

const (
	TeacherRequestTypeEdit   TeacherRequestType = "edit"
	TeacherRequestTypeDelete TeacherRequestType = "delete"
	TeacherRequestTypeOther  TeacherRequestType = "other"
)

type TeacherRequestStatus string

const (
	TeacherRequestStatusPending  TeacherRequestStatus = "pending"
	TeacherRequestStatusApproved TeacherRequestStatus = "approved"
	TeacherRequestStatusRejected TeacherRequestStatus = "rejected"
)

type TeacherRequest struct {
	bun.BaseModel `bun:"table:teacher_requests,alias:tr"`

	ID            uuid.UUID            `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	TeacherID     uuid.UUID            `bun:"teacher_id,notnull,type:uuid"`
	RequestType   TeacherRequestType   `bun:"request_type,notnull,type:teacher_request_type"`
	RequestData   map[string]any       `bun:"request_data,notnull,type:jsonb"`
	RequestReason *string              `bun:"request_reason"`
	Status        TeacherRequestStatus `bun:"status,notnull,type:teacher_request_status,default:pending"`
	ApprovedBy    *uuid.UUID           `bun:"approved_by,type:uuid"`
	ApprovedAt    *time.Time           `bun:"approved_at"`
	CreatedAt     time.Time            `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     time.Time            `bun:"updated_at,notnull,default:current_timestamp"`
}

type TeacherRequestUpdate struct {
	TeacherID     *uuid.UUID
	RequestType   *TeacherRequestType
	RequestData   *map[string]any
	RequestReason *string
	Status        *TeacherRequestStatus
	ApprovedBy    *uuid.UUID
	ApprovedAt    *time.Time
}

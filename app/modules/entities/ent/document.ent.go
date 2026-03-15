package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type DocumentStatus string

const (
	DocumentStatusPendingUpload DocumentStatus = "pending_upload"
	DocumentStatusActive        DocumentStatus = "active"
	DocumentStatusArchived      DocumentStatus = "archived"
)

type Document struct {
	bun.BaseModel `bun:"table:documents,alias:d"`

	ID                 uuid.UUID      `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID           uuid.UUID      `bun:"school_id,type:uuid,notnull"`
	StorageID          uuid.UUID      `bun:"storage_id,type:uuid,notnull"`
	StorageName        *string        `bun:"storage_name,scanonly"`
	OwnerMemberID      *uuid.UUID     `bun:"owner_member_id,type:uuid"`
	UploadedByMemberID uuid.UUID      `bun:"uploaded_by_member_id,type:uuid,notnull"`
	ObjectKey          string         `bun:"object_key,notnull"`
	FileName           string         `bun:"file_name,notnull"`
	ContentType        string         `bun:"content_type,notnull"`
	SizeBytes          int64          `bun:"size_bytes,notnull,default:0"`
	Status             DocumentStatus `bun:"status,notnull,default:'pending_upload'"`
	Metadata           map[string]any `bun:"metadata,type:jsonb"`
	CreatedAt          time.Time      `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt          time.Time      `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt          *time.Time     `bun:"deleted_at,soft_delete"`
}

type DocumentUpdate struct {
	OwnerMemberID *uuid.UUID
	FileName      *string
	ContentType   *string
	SizeBytes     *int64
	Status        *DocumentStatus
	Metadata      *map[string]any
}

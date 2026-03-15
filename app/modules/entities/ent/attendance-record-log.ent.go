package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type AttendanceRecordLog struct {
	bun.BaseModel `bun:"table:attendance_record_logs,alias:arl"`

	ID        uuid.UUID         `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	RecordID  uuid.UUID         `bun:"record_id,notnull,type:uuid"`
	OldStatus *AttendanceStatus `bun:"old_status,type:attendance_status"`
	NewStatus AttendanceStatus  `bun:"new_status,notnull,type:attendance_status"`
	ChangedBy *uuid.UUID        `bun:"changed_by,type:uuid"`
	ChangedAt time.Time         `bun:"changed_at,notnull,default:current_timestamp"`
	Reason    *string           `bun:"reason"`
	CreatedAt time.Time         `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time         `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt *time.Time        `bun:"deleted_at,soft_delete"`
}

type AttendanceRecordLogUpdate struct {
	RecordID  *uuid.UUID
	OldStatus *AttendanceStatus
	NewStatus *AttendanceStatus
	ChangedBy *uuid.UUID
	ChangedAt *time.Time
	Reason    *string
}

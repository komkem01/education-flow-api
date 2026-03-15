package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type AttendanceStatus string

const (
	AttendanceStatusPresent  AttendanceStatus = "present"
	AttendanceStatusLate     AttendanceStatus = "late"
	AttendanceStatusAbsent   AttendanceStatus = "absent"
	AttendanceStatusSick     AttendanceStatus = "sick"
	AttendanceStatusLeave    AttendanceStatus = "leave"
	AttendanceStatusActivity AttendanceStatus = "activity"
)

type AttendanceSource string

const (
	AttendanceSourceManual AttendanceSource = "manual"
	AttendanceSourceQR     AttendanceSource = "qr"
	AttendanceSourceRFID   AttendanceSource = "rfid"
	AttendanceSourceFace   AttendanceSource = "face"
	AttendanceSourceAPI    AttendanceSource = "api"
)

type AttendanceRecord struct {
	bun.BaseModel `bun:"table:attendance_records,alias:ar"`

	ID           uuid.UUID        `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SessionID    uuid.UUID        `bun:"session_id,notnull,type:uuid"`
	EnrollmentID uuid.UUID        `bun:"enrollment_id,notnull,type:uuid"`
	Status       AttendanceStatus `bun:"status,notnull,type:attendance_status,default:present"`
	Source       AttendanceSource `bun:"source,notnull,type:attendance_source,default:manual"`
	MarkedAt     time.Time        `bun:"marked_at,notnull,default:current_timestamp"`
	Remark       *string          `bun:"remark"`
	MarkedBy     *uuid.UUID       `bun:"marked_by,type:uuid"`
	CreatedAt    time.Time        `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt    time.Time        `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt    *time.Time       `bun:"deleted_at,soft_delete"`
}

type AttendanceRecordUpdate struct {
	SessionID    *uuid.UUID
	EnrollmentID *uuid.UUID
	Status       *AttendanceStatus
	Source       *AttendanceSource
	MarkedAt     *time.Time
	Remark       *string
	MarkedBy     *uuid.UUID
}

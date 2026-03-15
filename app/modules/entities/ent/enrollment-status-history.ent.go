package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type EnrollmentStatusHistory struct {
	bun.BaseModel `bun:"table:enrollment_status_histories,alias:esh"`

	ID           uuid.UUID                `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	EnrollmentID uuid.UUID                `bun:"enrollment_id,notnull,type:uuid"`
	FromStatus   *StudentEnrollmentStatus `bun:"from_status,type:student_enrollment_status"`
	ToStatus     StudentEnrollmentStatus  `bun:"to_status,notnull,type:student_enrollment_status"`
	ChangedAt    time.Time                `bun:"changed_at,notnull,default:current_timestamp"`
	ChangedBy    *uuid.UUID               `bun:"changed_by,type:uuid"`
	Reason       *string                  `bun:"reason"`
	CreatedAt    time.Time                `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt    time.Time                `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt    *time.Time               `bun:"deleted_at,soft_delete"`
}

type EnrollmentStatusHistoryUpdate struct {
	EnrollmentID *uuid.UUID
	FromStatus   *StudentEnrollmentStatus
	ToStatus     *StudentEnrollmentStatus
	ChangedAt    *time.Time
	ChangedBy    *uuid.UUID
	Reason       *string
}

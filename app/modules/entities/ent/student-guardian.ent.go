package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StudentGuardian struct {
	bun.BaseModel `bun:"table:student_guardians,alias:sgu"`

	ID                 uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	StudentID          uuid.UUID  `bun:"student_id,notnull,type:uuid"`
	GuardianID         uuid.UUID  `bun:"guardian_id,notnull,type:uuid"`
	Relationship       string     `bun:"relationship,notnull,type:guardian_relationship"`
	IsMainGuardian     bool       `bun:"is_main_guardian,notnull,default:false"`
	CanPickup          bool       `bun:"can_pickup,notnull,default:true"`
	IsEmergencyContact bool       `bun:"is_emergency_contact,notnull,default:false"`
	Note               *string    `bun:"note"`
	CreatedAt          time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt          time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt          *time.Time `bun:"deleted_at,soft_delete"`
}

type StudentGuardianUpdate struct {
	StudentID          *uuid.UUID
	GuardianID         *uuid.UUID
	Relationship       *string
	IsMainGuardian     *bool
	CanPickup          *bool
	IsEmergencyContact *bool
	Note               *string
}

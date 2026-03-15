package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeacherEmergencyContact struct {
	bun.BaseModel `bun:"table:teacher_emergency_contacts,alias:tec"`

	ID                   uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberTeacherID      uuid.UUID  `bun:"member_teacher_id,notnull,type:uuid"`
	EmergencyContactName string     `bun:"emergency_contact_name,notnull,type:varchar(255)"`
	Relationship         string     `bun:"relationship,notnull,type:varchar(100)"`
	PhonePrimary         string     `bun:"phone_primary,notnull,type:varchar(20)"`
	PhoneSecondary       *string    `bun:"phone_secondary,type:varchar(20)"`
	CanDecideMedical     bool       `bun:"can_decide_medical,notnull,default:false"`
	IsPrimary            bool       `bun:"is_primary,notnull,default:false"`
	CreatedAt            time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt            time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt            *time.Time `bun:"deleted_at,soft_delete"`
}

type TeacherEmergencyContactUpdate struct {
	MemberTeacherID      *uuid.UUID
	EmergencyContactName *string
	Relationship         *string
	PhonePrimary         *string
	PhoneSecondary       *string
	CanDecideMedical     *bool
	IsPrimary            *bool
}

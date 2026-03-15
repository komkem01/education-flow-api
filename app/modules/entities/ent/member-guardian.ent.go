package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberGuardian struct {
	bun.BaseModel `bun:"table:member_guardians,alias:mg"`

	ID          uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberID    *uuid.UUID `bun:"member_id,type:uuid"`
	SchoolID    uuid.UUID  `bun:"school_id,notnull,type:uuid"`
	GenderID    uuid.UUID  `bun:"gender_id,notnull,type:uuid"`
	PrefixID    uuid.UUID  `bun:"prefix_id,notnull,type:uuid"`
	FirstNameTH string     `bun:"first_name_th,notnull,type:varchar(255)"`
	LastNameTH  string     `bun:"last_name_th,notnull,type:varchar(255)"`
	FirstNameEN *string    `bun:"first_name_en,type:varchar(255)"`
	LastNameEN  *string    `bun:"last_name_en,type:varchar(255)"`
	CitizenID   *string    `bun:"citizen_id,type:varchar(13)"`
	Phone       *string    `bun:"phone,type:varchar(20)"`
	IsActive    bool       `bun:"is_active,notnull,default:true"`
	CreatedAt   time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt   *time.Time `bun:"deleted_at,soft_delete"`
}

type MemberGuardianUpdate struct {
	MemberID    *uuid.UUID
	SchoolID    *uuid.UUID
	GenderID    *uuid.UUID
	PrefixID    *uuid.UUID
	FirstNameTH *string
	LastNameTH  *string
	FirstNameEN *string
	LastNameEN  *string
	CitizenID   *string
	Phone       *string
	IsActive    *bool
}

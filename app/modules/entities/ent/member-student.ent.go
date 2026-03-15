package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberStudent struct {
	bun.BaseModel `bun:"table:member_students,alias:ms"`

	ID               uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberID         uuid.UUID  `bun:"member_id,notnull,type:uuid"`
	SchoolID         uuid.UUID  `bun:"school_id,notnull,type:uuid"`
	GenderID         uuid.UUID  `bun:"gender_id,notnull,type:uuid"`
	PrefixID         uuid.UUID  `bun:"prefix_id,notnull,type:uuid"`
	AdvisorTeacherID *uuid.UUID `bun:"advisor_teacher_id,type:uuid"`
	StudentCode      string     `bun:"student_code,notnull,type:varchar(50)"`
	FirstNameTH      string     `bun:"first_name_th,notnull,type:varchar(255)"`
	LastNameTH       string     `bun:"last_name_th,notnull,type:varchar(255)"`
	FirstNameEN      *string    `bun:"first_name_en,type:varchar(255)"`
	LastNameEN       *string    `bun:"last_name_en,type:varchar(255)"`
	CitizenID        *string    `bun:"citizen_id,type:varchar(13)"`
	Phone            *string    `bun:"phone,type:varchar(20)"`
	IsActive         bool       `bun:"is_active,notnull,default:true"`
	CreatedAt        time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt        time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt        *time.Time `bun:"deleted_at,soft_delete"`
}

type MemberStudentUpdate struct {
	MemberID         *uuid.UUID
	SchoolID         *uuid.UUID
	GenderID         *uuid.UUID
	PrefixID         *uuid.UUID
	AdvisorTeacherID *uuid.UUID
	StudentCode      *string
	FirstNameTH      *string
	LastNameTH       *string
	FirstNameEN      *string
	LastNameEN       *string
	CitizenID        *string
	Phone            *string
	IsActive         *bool
}

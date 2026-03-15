package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SubjectGroup struct {
	bun.BaseModel `bun:"table:subject_groups,alias:sg"`

	ID            uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID      uuid.UUID  `bun:"school_id,notnull,type:uuid"`
	Code          string     `bun:"code,notnull,type:varchar(30)"`
	NameTH        string     `bun:"name_th,notnull,type:varchar(255)"`
	NameEN        *string    `bun:"name_en,type:varchar(255)"`
	HeadTeacherID *uuid.UUID `bun:"head_teacher_id,type:uuid"`
	IsActive      bool       `bun:"is_active,notnull,default:true"`
	CreatedAt     time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt     *time.Time `bun:"deleted_at,soft_delete"`
}

type SubjectGroupUpdate struct {
	SchoolID      *uuid.UUID
	Code          *string
	NameTH        *string
	NameEN        *string
	HeadTeacherID *uuid.UUID
	IsActive      *bool
}

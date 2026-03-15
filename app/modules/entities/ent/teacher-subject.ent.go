package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeacherSubjectRole string

const (
	TeacherSubjectRolePrimary   TeacherSubjectRole = "primary"
	TeacherSubjectRoleAssistant TeacherSubjectRole = "assistant"
)

type TeacherSubject struct {
	bun.BaseModel `bun:"table:teacher_subjects,alias:ts"`

	ID        uuid.UUID          `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	TeacherID uuid.UUID          `bun:"teacher_id,notnull,type:uuid"`
	SubjectID uuid.UUID          `bun:"subject_id,notnull,type:uuid"`
	Role      TeacherSubjectRole `bun:"role,notnull,type:teacher_subject_role,default:primary"`
	IsActive  bool               `bun:"is_active,notnull,default:true"`
	CreatedAt time.Time          `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time          `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt *time.Time         `bun:"deleted_at,soft_delete"`
}

type TeacherSubjectUpdate struct {
	TeacherID *uuid.UUID
	SubjectID *uuid.UUID
	Role      *TeacherSubjectRole
	IsActive  *bool
}

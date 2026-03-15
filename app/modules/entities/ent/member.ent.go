package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberRole string

const (
	MemberRoleSuperadmin MemberRole = "superadmin"
	MemberRoleAdmin      MemberRole = "admin"
	MemberRoleStaff      MemberRole = "staff"
	MemberRoleTeacher    MemberRole = "teacher"
	MemberRoleStudent    MemberRole = "student"
)

type Member struct {
	bun.BaseModel `bun:"table:members,alias:m"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID  uuid.UUID  `bun:"school_id,type:uuid,notnull"`
	Email     string     `bun:"email,notnull,unique"`
	Password  string     `bun:"password,notnull"`
	Role      MemberRole `bun:"role,notnull,default:'admin'"`
	IsActive  bool       `bun:"is_active,notnull,default:false"`
	LastLogin *time.Time `bun:"last_login"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
}

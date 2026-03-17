package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SchoolDepartment struct {
	bun.BaseModel `bun:"table:school_departments,alias:sd"`

	ID           uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID     uuid.UUID  `bun:"school_id,notnull,type:uuid"`
	DepartmentID uuid.UUID  `bun:"department_id,notnull,type:uuid"`
	Code         string     `bun:"code,notnull,type:varchar(50)"`
	CustomName   *string    `bun:"custom_name,type:varchar(255)"`
	IsActive     bool       `bun:"is_active,notnull,default:true"`
	CreatedAt    time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt    time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt    *time.Time `bun:"deleted_at,soft_delete"`
}

type SchoolDepartmentUpdate struct {
	SchoolID     *uuid.UUID
	DepartmentID *uuid.UUID
	Code         *string
	CustomName   *string
	IsActive     *bool
}

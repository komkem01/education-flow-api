package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeacherExperience struct {
	bun.BaseModel `bun:"table:teacher_experiences,alias:tx"`

	ID               uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	TeacherID        uuid.UUID  `bun:"teacher_id,notnull,type:uuid"`
	SchoolName       string     `bun:"school_name,notnull,type:varchar(255)"`
	Position         string     `bun:"position,notnull,type:varchar(255)"`
	DepartmentName   *string    `bun:"department_name,type:varchar(255)"`
	StartDate        time.Time  `bun:"start_date,notnull,type:date"`
	EndDate          *time.Time `bun:"end_date,type:date"`
	IsCurrent        bool       `bun:"is_current,notnull,default:false"`
	Responsibilities *string    `bun:"responsibilities"`
	Achievements     *string    `bun:"achievements"`
	SortOrder        int        `bun:"sort_order,notnull,default:0"`
	IsActive         bool       `bun:"is_active,notnull,default:true"`
	CreatedAt        time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt        time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt        *time.Time `bun:"deleted_at,soft_delete"`
}

type TeacherExperienceUpdate struct {
	TeacherID        *uuid.UUID
	SchoolName       *string
	Position         *string
	DepartmentName   *string
	StartDate        *time.Time
	EndDate          *time.Time
	IsCurrent        *bool
	Responsibilities *string
	Achievements     *string
	SortOrder        *int
	IsActive         *bool
}

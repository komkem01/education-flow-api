package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Subject struct {
	bun.BaseModel `bun:"table:subjects,alias:su"`

	ID             uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID       uuid.UUID  `bun:"school_id,notnull,type:uuid"`
	SubjectGroupID uuid.UUID  `bun:"subject_group_id,notnull,type:uuid"`
	Code           string     `bun:"code,notnull,type:varchar(30)"`
	NameTH         string     `bun:"name_th,notnull,type:varchar(255)"`
	NameEN         *string    `bun:"name_en,type:varchar(255)"`
	Credit         float64    `bun:"credit,notnull"`
	HoursPerWeek   *int       `bun:"hours_per_week"`
	IsElective     bool       `bun:"is_elective,notnull,default:false"`
	IsActive       bool       `bun:"is_active,notnull,default:true"`
	CreatedAt      time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt      time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt      *time.Time `bun:"deleted_at,soft_delete"`
}

type SubjectUpdate struct {
	SchoolID       *uuid.UUID
	SubjectGroupID *uuid.UUID
	Code           *string
	NameTH         *string
	NameEN         *string
	Credit         *float64
	HoursPerWeek   *int
	IsElective     *bool
	IsActive       *bool
}

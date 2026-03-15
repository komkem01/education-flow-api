package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type AcademicYear struct {
	bun.BaseModel `bun:"table:academic_years,alias:ay"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID  uuid.UUID  `bun:"school_id,type:uuid,notnull"`
	Year      string     `bun:"year,notnull"`
	StartDate time.Time  `bun:"start_date,type:date,notnull"`
	EndDate   time.Time  `bun:"end_date,type:date,notnull"`
	IsActive  bool       `bun:"is_active,notnull,default:false"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
}

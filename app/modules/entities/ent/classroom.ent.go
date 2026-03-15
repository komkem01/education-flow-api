package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Classroom struct {
	bun.BaseModel `bun:"table:classrooms,alias:c"`

	ID                uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID          uuid.UUID  `bun:"school_id,notnull,type:uuid"`
	AcademicYearID    uuid.UUID  `bun:"academic_year_id,notnull,type:uuid"`
	Level             string     `bun:"level,notnull,type:varchar(50)"`
	RoomNo            string     `bun:"room_no,notnull,type:varchar(20)"`
	Name              string     `bun:"name,notnull,type:varchar(120)"`
	HomeroomTeacherID *uuid.UUID `bun:"homeroom_teacher_id,type:uuid"`
	Capacity          *int       `bun:"capacity"`
	IsActive          bool       `bun:"is_active,notnull,default:true"`
	CreatedAt         time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt         time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt         *time.Time `bun:"deleted_at,soft_delete"`
}

type ClassroomUpdate struct {
	SchoolID          *uuid.UUID
	AcademicYearID    *uuid.UUID
	Level             *string
	RoomNo            *string
	Name              *string
	HomeroomTeacherID *uuid.UUID
	Capacity          *int
	IsActive          *bool
}

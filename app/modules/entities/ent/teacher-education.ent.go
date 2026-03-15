package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeacherDegree string

const (
	TeacherDegreeBachelor TeacherDegree = "ปริญญาตรี"
	TeacherDegreeMaster   TeacherDegree = "ปริญญาโท"
	TeacherDegreeDoctor   TeacherDegree = "ปริญญาเอก"
)

type TeacherEducation struct {
	bun.BaseModel `bun:"table:teacher_educations,alias:te"`

	ID             uuid.UUID     `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	TeacherID      uuid.UUID     `bun:"teacher_id,notnull,type:uuid"`
	Degree         TeacherDegree `bun:"degree,notnull,type:teacher_degree"`
	Major          string        `bun:"major,notnull"`
	University     string        `bun:"university,notnull"`
	GraduationYear string        `bun:"graduation_year,notnull"`
	CreatedAt      time.Time     `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt      time.Time     `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt      *time.Time    `bun:"deleted_at,soft_delete"`
}

type TeacherEducationUpdate struct {
	TeacherID      *uuid.UUID
	Degree         *TeacherDegree
	Major          *string
	University     *string
	GraduationYear *string
}

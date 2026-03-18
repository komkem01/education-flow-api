package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type EnrollmentSubject struct {
	bun.BaseModel `bun:"table:enrollment_subjects,alias:es"`

	ID            uuid.UUID               `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	EnrollmentID  uuid.UUID               `bun:"enrollment_id,notnull,type:uuid"`
	SubjectID     uuid.UUID               `bun:"subject_id,notnull,type:uuid"`
	TeacherID     *uuid.UUID              `bun:"teacher_id,type:uuid"`
	IsPrimary     bool                    `bun:"is_primary,notnull,default:false"`
	Status        StudentEnrollmentStatus `bun:"status,notnull,type:student_enrollment_status,default:active"`
	MidtermScore  *float64                `bun:"midterm_score,type:numeric(5,2)"`
	FinalScore    *float64                `bun:"final_score,type:numeric(5,2)"`
	ActivityScore *float64                `bun:"activity_score,type:numeric(5,2)"`
	TotalScore    *float64                `bun:"-"`
	GradeNumeric  *string                 `bun:"-"`
	GradeLetter   *string                 `bun:"-"`
	CreatedAt     time.Time               `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     time.Time               `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt     *time.Time              `bun:"deleted_at,soft_delete"`
}

type EnrollmentSubjectUpdate struct {
	EnrollmentID  *uuid.UUID
	SubjectID     *uuid.UUID
	TeacherID     *uuid.UUID
	IsPrimary     *bool
	Status        *StudentEnrollmentStatus
	MidtermScore  *float64
	FinalScore    *float64
	ActivityScore *float64
}

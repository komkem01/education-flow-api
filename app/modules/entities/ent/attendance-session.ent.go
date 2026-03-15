package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type AttendanceMode string

const (
	AttendanceModeHomeroom AttendanceMode = "homeroom"
	AttendanceModeSubject  AttendanceMode = "subject"
	AttendanceModeActivity AttendanceMode = "activity"
)

type AttendanceSession struct {
	bun.BaseModel `bun:"table:attendance_sessions,alias:as"`

	ID             uuid.UUID      `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID       uuid.UUID      `bun:"school_id,notnull,type:uuid"`
	AcademicYearID uuid.UUID      `bun:"academic_year_id,notnull,type:uuid"`
	ClassroomID    uuid.UUID      `bun:"classroom_id,notnull,type:uuid"`
	SubjectID      *uuid.UUID     `bun:"subject_id,type:uuid"`
	TeacherID      *uuid.UUID     `bun:"teacher_id,type:uuid"`
	SessionDate    time.Time      `bun:"session_date,notnull,type:date"`
	PeriodNo       int            `bun:"period_no,notnull"`
	Mode           AttendanceMode `bun:"mode,notnull,type:attendance_mode,default:homeroom"`
	StartedAt      time.Time      `bun:"started_at,notnull,default:current_timestamp"`
	ClosedAt       *time.Time     `bun:"closed_at"`
	Note           *string        `bun:"note"`
	CreatedAt      time.Time      `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt      time.Time      `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt      *time.Time     `bun:"deleted_at,soft_delete"`
}

type AttendanceSessionUpdate struct {
	SchoolID       *uuid.UUID
	AcademicYearID *uuid.UUID
	ClassroomID    *uuid.UUID
	SubjectID      *uuid.UUID
	TeacherID      *uuid.UUID
	SessionDate    *time.Time
	PeriodNo       *int
	Mode           *AttendanceMode
	StartedAt      *time.Time
	ClosedAt       *time.Time
	Note           *string
}

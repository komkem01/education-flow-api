package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StudentEnrollmentStatus string

const (
	StudentEnrollmentStatusActive      StudentEnrollmentStatus = "active"
	StudentEnrollmentStatusTransferred StudentEnrollmentStatus = "transferred"
	StudentEnrollmentStatusGraduated   StudentEnrollmentStatus = "graduated"
	StudentEnrollmentStatusDropped     StudentEnrollmentStatus = "dropped"
)

type EnrollmentType string

const (
	EnrollmentTypeNew        EnrollmentType = "new"
	EnrollmentTypeTransferIn EnrollmentType = "transfer_in"
	EnrollmentTypeRepeat     EnrollmentType = "repeat"
	EnrollmentTypeReturn     EnrollmentType = "return"
)

type EnrollmentExitReason string

const (
	EnrollmentExitReasonTransferOut EnrollmentExitReason = "transfer_out"
	EnrollmentExitReasonGraduated   EnrollmentExitReason = "graduated"
	EnrollmentExitReasonDropped     EnrollmentExitReason = "dropped"
	EnrollmentExitReasonLeave       EnrollmentExitReason = "leave"
)

type StudentEnrollment struct {
	bun.BaseModel `bun:"table:student_enrollments,alias:se"`

	ID                   uuid.UUID               `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	StudentID            uuid.UUID               `bun:"student_id,notnull,type:uuid"`
	SchoolID             uuid.UUID               `bun:"school_id,notnull,type:uuid"`
	AcademicYearID       uuid.UUID               `bun:"academic_year_id,notnull,type:uuid"`
	ClassroomID          uuid.UUID               `bun:"classroom_id,notnull,type:uuid"`
	EnrolledAt           *time.Time              `bun:"enrolled_at,type:date"`
	ExitedAt             *time.Time              `bun:"exited_at,type:date"`
	Status               StudentEnrollmentStatus `bun:"status,notnull,type:student_enrollment_status,default:active"`
	EnrollmentType       EnrollmentType          `bun:"enrollment_type,notnull,type:enrollment_type,default:new"`
	ExitReason           *EnrollmentExitReason   `bun:"exit_reason,type:enrollment_exit_reason"`
	ExitNote             *string                 `bun:"exit_note"`
	PreviousEnrollmentID *uuid.UUID              `bun:"previous_enrollment_id,type:uuid"`
	RollNo               *string                 `bun:"roll_no,type:varchar(10)"`
	ApprovedBy           *uuid.UUID              `bun:"approved_by,type:uuid"`
	ApprovedAt           *time.Time              `bun:"approved_at"`
	ApprovalNote         *string                 `bun:"approval_note"`
	CreatedBy            *uuid.UUID              `bun:"created_by,type:uuid"`
	UpdatedBy            *uuid.UUID              `bun:"updated_by,type:uuid"`
	CreatedAt            time.Time               `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt            time.Time               `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt            *time.Time              `bun:"deleted_at,soft_delete"`
}

type StudentEnrollmentUpdate struct {
	StudentID            *uuid.UUID
	SchoolID             *uuid.UUID
	AcademicYearID       *uuid.UUID
	ClassroomID          *uuid.UUID
	EnrolledAt           *time.Time
	ExitedAt             *time.Time
	Status               *StudentEnrollmentStatus
	EnrollmentType       *EnrollmentType
	ExitReason           *EnrollmentExitReason
	ExitNote             *string
	PreviousEnrollmentID *uuid.UUID
	RollNo               *string
	ApprovedBy           *uuid.UUID
	ApprovedAt           *time.Time
	ApprovalNote         *string
	CreatedBy            *uuid.UUID
	UpdatedBy            *uuid.UUID
}

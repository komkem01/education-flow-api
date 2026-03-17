package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberManagement struct {
	bun.BaseModel `bun:"table:member_managements,alias:mm"`

	ID                 uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberID           uuid.UUID  `bun:"member_id,notnull,type:uuid"`
	EmployeeCode       string     `bun:"employee_code,notnull,type:varchar(50)"`
	GenderID           *uuid.UUID `bun:"gender_id,type:uuid"`
	PrefixID           *uuid.UUID `bun:"prefix_id,type:uuid"`
	FirstName          *string    `bun:"first_name,type:varchar(255)"`
	LastName           *string    `bun:"last_name,type:varchar(255)"`
	Phone              *string    `bun:"phone,type:varchar(20)"`
	Position           string     `bun:"position,notnull,type:varchar(255)"`
	StartWorkDate      time.Time  `bun:"start_work_date,notnull,type:date"`
	SchoolDepartmentID uuid.UUID  `bun:"school_department_id,notnull,type:uuid"`
	DepartmentID       uuid.UUID  `bun:"department_id,notnull,type:uuid"`
	IsActive           bool       `bun:"is_active,notnull,default:true"`
	CreatedAt          time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt          time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt          *time.Time `bun:"deleted_at,soft_delete"`
}

type MemberManagementUpdate struct {
	MemberID           *uuid.UUID
	EmployeeCode       *string
	GenderID           *uuid.UUID
	PrefixID           *uuid.UUID
	FirstName          *string
	LastName           *string
	Phone              *string
	Position           *string
	StartWorkDate      *time.Time
	SchoolDepartmentID *uuid.UUID
	DepartmentID       *uuid.UUID
	IsActive           *bool
}

type ManagementRegistrationInput struct {
	MemberEmail        string
	MemberPasswordHash string
	MemberSchoolID     uuid.UUID
	MemberRole         MemberRole
	MemberIsActive     bool
	MemberLastLogin    *time.Time

	ManagementEmployeeCode       string
	ManagementGenderID           *uuid.UUID
	ManagementPrefixID           *uuid.UUID
	ManagementFirstName          *string
	ManagementLastName           *string
	ManagementPhone              *string
	ManagementPosition           string
	ManagementStartWorkDate      time.Time
	ManagementSchoolDepartmentID uuid.UUID
	ManagementDepartmentID       uuid.UUID
	ManagementIsActive           bool

	RequestedBy     uuid.UUID
	RequestedByRole ApprovalActorRole
	RequestReason   *string
}

type ManagementRegistrationResult struct {
	Member         *Member
	Management     *MemberManagement
	Approval       *ApprovalRequest
	ApprovalAction *ApprovalAction
}

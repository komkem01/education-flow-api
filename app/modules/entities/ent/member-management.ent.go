package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberManagement struct {
	bun.BaseModel `bun:"table:member_managements,alias:mm"`

	ID            uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberID      uuid.UUID  `bun:"member_id,notnull,type:uuid"`
	EmployeeCode  string     `bun:"employee_code,notnull,type:varchar(50)"`
	Position      string     `bun:"position,notnull,type:varchar(255)"`
	StartWorkDate time.Time  `bun:"start_work_date,notnull,type:date"`
	DepartmentID  uuid.UUID  `bun:"department_id,notnull,type:uuid"`
	IsActive      bool       `bun:"is_active,notnull,default:true"`
	CreatedAt     time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt     *time.Time `bun:"deleted_at,soft_delete"`
}

type MemberManagementUpdate struct {
	MemberID      *uuid.UUID
	EmployeeCode  *string
	Position      *string
	StartWorkDate *time.Time
	DepartmentID  *uuid.UUID
	IsActive      *bool
}

type ManagementRegistrationInput struct {
	MemberEmail        string
	MemberPasswordHash string
	MemberSchoolID     uuid.UUID
	MemberRole         MemberRole
	MemberIsActive     bool
	MemberLastLogin    *time.Time

	ManagementEmployeeCode  string
	ManagementPosition      string
	ManagementStartWorkDate time.Time
	ManagementDepartmentID  uuid.UUID
	ManagementIsActive      bool

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

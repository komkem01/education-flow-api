package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberTeacher struct {
	bun.BaseModel `bun:"table:member_teachers,alias:mt"`

	ID               uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberID         uuid.UUID  `bun:"member_id,notnull,type:uuid"`
	GenderID         uuid.UUID  `bun:"gender_id,notnull,type:uuid"`
	PrefixID         uuid.UUID  `bun:"prefix_id,notnull,type:uuid"`
	Code             string     `bun:"code,notnull,type:varchar(255)"`
	CitizenID        string     `bun:"citizen_id,notnull,type:varchar(13)"`
	FirstNameTH      string     `bun:"first_name_th,notnull,type:varchar(255)"`
	LastNameTH       string     `bun:"last_name_th,notnull,type:varchar(255)"`
	FirstNameEN      string     `bun:"first_name_en,notnull,type:varchar(255)"`
	LastNameEN       string     `bun:"last_name_en,notnull,type:varchar(255)"`
	Phone            string     `bun:"phone,notnull,type:varchar(20)"`
	Position         string     `bun:"position,notnull,type:varchar(255)"`
	AcademicStanding string     `bun:"academic_standing,notnull,type:varchar(255)"`
	DepartmentID     uuid.UUID  `bun:"department,notnull,type:uuid"`
	StartDate        time.Time  `bun:"start_date,notnull,type:date"`
	EndDate          *time.Time `bun:"end_date,type:date"`
	IsActive         bool       `bun:"is_active,notnull,default:true"`
	CreatedAt        time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt        time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt        *time.Time `bun:"deleted_at,soft_delete"`
}

type MemberTeacherUpdate struct {
	MemberID         *uuid.UUID
	GenderID         *uuid.UUID
	PrefixID         *uuid.UUID
	Code             *string
	CitizenID        *string
	FirstNameTH      *string
	LastNameTH       *string
	FirstNameEN      *string
	LastNameEN       *string
	Phone            *string
	Position         *string
	AcademicStanding *string
	DepartmentID     *uuid.UUID
	StartDate        *time.Time
	EndDate          *time.Time
	IsActive         *bool
}

type TeacherRegistrationInput struct {
	MemberEmail        string
	MemberPasswordHash string
	MemberSchoolID     uuid.UUID
	MemberRole         MemberRole
	MemberIsActive     bool
	MemberLastLogin    *time.Time

	TeacherGenderID         uuid.UUID
	TeacherPrefixID         uuid.UUID
	TeacherCode             string
	TeacherCitizenID        string
	TeacherFirstNameTH      string
	TeacherLastNameTH       string
	TeacherFirstNameEN      string
	TeacherLastNameEN       string
	TeacherPhone            string
	TeacherPosition         string
	TeacherAcademicStanding string
	TeacherDepartmentID     uuid.UUID
	TeacherStartDate        time.Time
	TeacherEndDate          *time.Time
	TeacherIsActive         bool
}

type TeacherRegistrationResult struct {
	Member  *Member
	Teacher *MemberTeacher
}

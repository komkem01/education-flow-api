package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeacherLicenseStatus string

const (
	TeacherLicenseStatusActive    TeacherLicenseStatus = "active"
	TeacherLicenseStatusSuspended TeacherLicenseStatus = "suspended"
	TeacherLicenseStatusExpired   TeacherLicenseStatus = "expired"
	TeacherLicenseStatusRevoked   TeacherLicenseStatus = "revoked"
)

type TeacherLicense struct {
	bun.BaseModel `bun:"table:teacher_licenses,alias:tl"`

	ID            uuid.UUID            `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	TeacherID     uuid.UUID            `bun:"teacher_id,notnull,type:uuid"`
	LicenseNo     string               `bun:"license_no,notnull,type:varchar(100)"`
	IssuedAt      *time.Time           `bun:"issued_at,type:date"`
	ExpiresAt     *time.Time           `bun:"expires_at,type:date"`
	LicenseStatus TeacherLicenseStatus `bun:"license_status,notnull,type:teacher_license_status,default:active"`
	IssuedBy      *string              `bun:"issued_by,type:varchar(255)"`
	Note          *string              `bun:"note"`
	CreatedAt     time.Time            `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     time.Time            `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt     *time.Time           `bun:"deleted_at,soft_delete"`
}

type TeacherLicenseUpdate struct {
	TeacherID     *uuid.UUID
	LicenseNo     *string
	IssuedAt      *time.Time
	ExpiresAt     *time.Time
	LicenseStatus *TeacherLicenseStatus
	IssuedBy      *string
	Note          *string
}

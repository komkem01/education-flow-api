package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StudentProfile struct {
	bun.BaseModel `bun:"table:student_profiles,alias:sp"`

	ID                    uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	StudentID             uuid.UUID  `bun:"student_id,notnull,type:uuid"`
	BirthDate             *time.Time `bun:"birth_date,type:date"`
	Nationality           *string    `bun:"nationality,type:varchar(100)"`
	Religion              *string    `bun:"religion,type:varchar(100)"`
	AddressCurrent        *string    `bun:"address_current"`
	AddressRegistered     *string    `bun:"address_registered"`
	EmergencyContactName  *string    `bun:"emergency_contact_name,type:varchar(255)"`
	EmergencyContactPhone *string    `bun:"emergency_contact_phone,type:varchar(20)"`
	CreatedAt             time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt             time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt             *time.Time `bun:"deleted_at,soft_delete"`
}

type StudentProfileUpdate struct {
	StudentID             *uuid.UUID
	BirthDate             *time.Time
	Nationality           *string
	Religion              *string
	AddressCurrent        *string
	AddressRegistered     *string
	EmergencyContactName  *string
	EmergencyContactPhone *string
}

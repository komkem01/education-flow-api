package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StudentHealthProfile struct {
	bun.BaseModel `bun:"table:student_health_profiles,alias:shp"`

	ID             uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	StudentID      uuid.UUID  `bun:"student_id,notnull,type:uuid"`
	BloodType      *string    `bun:"blood_type,type:blood_type"`
	AllergyInfo    *string    `bun:"allergy_info"`
	ChronicDisease *string    `bun:"chronic_disease"`
	MedicalNote    *string    `bun:"medical_note"`
	CreatedAt      time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt      time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt      *time.Time `bun:"deleted_at,soft_delete"`
}

type StudentHealthProfileUpdate struct {
	StudentID      *uuid.UUID
	BloodType      *string
	AllergyInfo    *string
	ChronicDisease *string
	MedicalNote    *string
}

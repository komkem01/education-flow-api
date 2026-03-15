package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TeacherHealthProfile struct {
	bun.BaseModel `bun:"table:teacher_health_profiles,alias:thp"`

	ID                 uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberTeacherID    uuid.UUID  `bun:"member_teacher_id,notnull,type:uuid"`
	BloodType          *string    `bun:"blood_type,type:blood_type"`
	AllergyInfo        *string    `bun:"allergy_info"`
	ChronicDisease     *string    `bun:"chronic_disease"`
	MedicationNote     *string    `bun:"medication_note"`
	FitnessForWorkNote *string    `bun:"fitness_for_work_note"`
	CreatedAt          time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt          time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt          *time.Time `bun:"deleted_at,soft_delete"`
}

type TeacherHealthProfileUpdate struct {
	MemberTeacherID    *uuid.UUID
	BloodType          *string
	AllergyInfo        *string
	ChronicDisease     *string
	MedicationNote     *string
	FitnessForWorkNote *string
}

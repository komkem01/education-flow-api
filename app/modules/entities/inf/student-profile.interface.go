package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type StudentProfileEntity interface {
	CreateStudentProfile(ctx context.Context, data *ent.StudentProfile) (*ent.StudentProfile, error)
	GetStudentProfileByID(ctx context.Context, id uuid.UUID) (*ent.StudentProfile, error)
	ListStudentProfiles(ctx context.Context, req *base.RequestPaginate, studentID *uuid.UUID) ([]*ent.StudentProfile, *base.ResponsePaginate, error)
	UpdateStudentProfileByID(ctx context.Context, id uuid.UUID, data *ent.StudentProfileUpdate) (*ent.StudentProfile, error)
	SoftDeleteStudentProfileByID(ctx context.Context, id uuid.UUID) error
}

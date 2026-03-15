package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type StudentHealthProfileEntity interface {
	CreateStudentHealthProfile(ctx context.Context, data *ent.StudentHealthProfile) (*ent.StudentHealthProfile, error)
	GetStudentHealthProfileByID(ctx context.Context, id uuid.UUID) (*ent.StudentHealthProfile, error)
	ListStudentHealthProfiles(ctx context.Context, req *base.RequestPaginate, studentID *uuid.UUID, bloodType *string) ([]*ent.StudentHealthProfile, *base.ResponsePaginate, error)
	UpdateStudentHealthProfileByID(ctx context.Context, id uuid.UUID, data *ent.StudentHealthProfileUpdate) (*ent.StudentHealthProfile, error)
	SoftDeleteStudentHealthProfileByID(ctx context.Context, id uuid.UUID) error
}

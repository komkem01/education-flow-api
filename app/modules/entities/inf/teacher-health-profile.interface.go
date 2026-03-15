package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type TeacherHealthProfileEntity interface {
	CreateTeacherHealthProfile(ctx context.Context, data *ent.TeacherHealthProfile) (*ent.TeacherHealthProfile, error)
	GetTeacherHealthProfileByID(ctx context.Context, id uuid.UUID) (*ent.TeacherHealthProfile, error)
	ListTeacherHealthProfiles(ctx context.Context, req *base.RequestPaginate, memberTeacherID *uuid.UUID, bloodType *string) ([]*ent.TeacherHealthProfile, *base.ResponsePaginate, error)
	UpdateTeacherHealthProfileByID(ctx context.Context, id uuid.UUID, data *ent.TeacherHealthProfileUpdate) (*ent.TeacherHealthProfile, error)
	SoftDeleteTeacherHealthProfileByID(ctx context.Context, id uuid.UUID) error
}

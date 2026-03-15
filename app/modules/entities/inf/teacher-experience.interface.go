package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type TeacherExperienceEntity interface {
	CreateTeacherExperience(ctx context.Context, data *ent.TeacherExperience) (*ent.TeacherExperience, error)
	GetTeacherExperienceByID(ctx context.Context, id uuid.UUID) (*ent.TeacherExperience, error)
	ListTeacherExperiences(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, isCurrent *bool, isActive *bool) ([]*ent.TeacherExperience, *base.ResponsePaginate, error)
	UpdateTeacherExperienceByID(ctx context.Context, id uuid.UUID, data *ent.TeacherExperienceUpdate) (*ent.TeacherExperience, error)
	SoftDeleteTeacherExperienceByID(ctx context.Context, id uuid.UUID) error
}

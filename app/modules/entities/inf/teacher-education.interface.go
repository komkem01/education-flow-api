package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type TeacherEducationEntity interface {
	CreateTeacherEducation(ctx context.Context, data *ent.TeacherEducation) (*ent.TeacherEducation, error)
	GetTeacherEducationByID(ctx context.Context, id uuid.UUID) (*ent.TeacherEducation, error)
	ListTeacherEducations(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, degree *ent.TeacherDegree) ([]*ent.TeacherEducation, *base.ResponsePaginate, error)
	UpdateTeacherEducationByID(ctx context.Context, id uuid.UUID, data *ent.TeacherEducationUpdate) (*ent.TeacherEducation, error)
	SoftDeleteTeacherEducationByID(ctx context.Context, id uuid.UUID) error
}

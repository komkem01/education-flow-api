package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type SubjectGroupEntity interface {
	CreateSubjectGroup(ctx context.Context, data *ent.SubjectGroup) (*ent.SubjectGroup, error)
	GetSubjectGroupByID(ctx context.Context, id uuid.UUID) (*ent.SubjectGroup, error)
	ListSubjectGroups(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, headTeacherID *uuid.UUID) ([]*ent.SubjectGroup, *base.ResponsePaginate, error)
	UpdateSubjectGroupByID(ctx context.Context, id uuid.UUID, data *ent.SubjectGroupUpdate) (*ent.SubjectGroup, error)
	SoftDeleteSubjectGroupByID(ctx context.Context, id uuid.UUID) error
}

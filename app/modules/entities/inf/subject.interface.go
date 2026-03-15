package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type SubjectEntity interface {
	CreateSubject(ctx context.Context, data *ent.Subject) (*ent.Subject, error)
	GetSubjectByID(ctx context.Context, id uuid.UUID) (*ent.Subject, error)
	ListSubjects(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, subjectGroupID *uuid.UUID, isElective *bool) ([]*ent.Subject, *base.ResponsePaginate, error)
	UpdateSubjectByID(ctx context.Context, id uuid.UUID, data *ent.SubjectUpdate) (*ent.Subject, error)
	SoftDeleteSubjectByID(ctx context.Context, id uuid.UUID) error
}

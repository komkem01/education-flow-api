package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type PrefixEntity interface {
	CreatePrefix(ctx context.Context, genderID uuid.UUID, name string, isActive bool) (*ent.Prefix, error)
	GetPrefixByID(ctx context.Context, id uuid.UUID) (*ent.Prefix, error)
	ListPrefixes(ctx context.Context, req *base.RequestPaginate, isActive *bool, genderID *uuid.UUID) ([]*ent.Prefix, *base.ResponsePaginate, error)
	UpdatePrefixByID(ctx context.Context, id uuid.UUID, genderID *uuid.UUID, name *string, isActive *bool) (*ent.Prefix, error)
	SoftDeletePrefixByID(ctx context.Context, id uuid.UUID) error
}

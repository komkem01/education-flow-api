package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type GenderEntity interface {
	CreateGender(ctx context.Context, name string, isActive bool) (*ent.Gender, error)
	GetGenderByID(ctx context.Context, id uuid.UUID) (*ent.Gender, error)
	ListGenders(ctx context.Context, req *base.RequestPaginate, isActive *bool) ([]*ent.Gender, *base.ResponsePaginate, error)
	UpdateGenderByID(ctx context.Context, id uuid.UUID, name *string, isActive *bool) (*ent.Gender, error)
	SoftDeleteGenderByID(ctx context.Context, id uuid.UUID) error
}

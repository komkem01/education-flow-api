package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type SchoolEntity interface {
	CreateSchool(ctx context.Context, name string, logoURL *string, themeColor *string, address *string, description *string) (*ent.School, error)
	GetSchoolByID(ctx context.Context, id uuid.UUID) (*ent.School, error)
	ListSchools(ctx context.Context, req *base.RequestPaginate) ([]*ent.School, *base.ResponsePaginate, error)
	UpdateSchoolByID(ctx context.Context, id uuid.UUID, name *string, logoURL *string, themeColor *string, address *string, description *string) (*ent.School, error)
	SoftDeleteSchoolByID(ctx context.Context, id uuid.UUID) error
}

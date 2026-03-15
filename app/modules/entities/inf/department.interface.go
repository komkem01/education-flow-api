package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type DepartmentEntity interface {
	CreateDepartment(ctx context.Context, name string, isActive bool) (*ent.Department, error)
	GetDepartmentByID(ctx context.Context, id uuid.UUID) (*ent.Department, error)
	ListDepartments(ctx context.Context, req *base.RequestPaginate, isActive *bool) ([]*ent.Department, *base.ResponsePaginate, error)
	UpdateDepartmentByID(ctx context.Context, id uuid.UUID, name *string, isActive *bool) (*ent.Department, error)
	SoftDeleteDepartmentByID(ctx context.Context, id uuid.UUID) error
}

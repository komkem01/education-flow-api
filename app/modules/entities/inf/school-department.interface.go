package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type SchoolDepartmentEntity interface {
	CreateSchoolDepartment(ctx context.Context, data *ent.SchoolDepartment) (*ent.SchoolDepartment, error)
	GetSchoolDepartmentByID(ctx context.Context, id uuid.UUID) (*ent.SchoolDepartment, error)
	GetSchoolDepartmentBySchoolAndDepartment(ctx context.Context, schoolID uuid.UUID, departmentID uuid.UUID) (*ent.SchoolDepartment, error)
	ListSchoolDepartments(ctx context.Context, req *base.RequestPaginate, schoolID *uuid.UUID, departmentID *uuid.UUID, isActive *bool) ([]*ent.SchoolDepartment, *base.ResponsePaginate, error)
	UpdateSchoolDepartmentByID(ctx context.Context, id uuid.UUID, data *ent.SchoolDepartmentUpdate) (*ent.SchoolDepartment, error)
	SoftDeleteSchoolDepartmentByID(ctx context.Context, id uuid.UUID) error
}

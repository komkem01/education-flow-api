package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type MemberManagementEntity interface {
	CreateMemberManagement(ctx context.Context, data *ent.MemberManagement) (*ent.MemberManagement, error)
	RegisterManagement(ctx context.Context, data *ent.ManagementRegistrationInput) (*ent.ManagementRegistrationResult, error)
	GetMemberManagementByID(ctx context.Context, id uuid.UUID) (*ent.MemberManagement, error)
	ListMemberManagements(ctx context.Context, req *base.RequestPaginate, isActive *bool, memberID *uuid.UUID, departmentID *uuid.UUID, schoolDepartmentID *uuid.UUID) ([]*ent.MemberManagement, *base.ResponsePaginate, error)
	UpdateMemberManagementByID(ctx context.Context, id uuid.UUID, data *ent.MemberManagementUpdate) (*ent.MemberManagement, error)
	SoftDeleteMemberManagementByID(ctx context.Context, id uuid.UUID) error
}

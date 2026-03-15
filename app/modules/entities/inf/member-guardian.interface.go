package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type MemberGuardianEntity interface {
	CreateMemberGuardian(ctx context.Context, data *ent.MemberGuardian) (*ent.MemberGuardian, error)
	GetMemberGuardianByID(ctx context.Context, id uuid.UUID) (*ent.MemberGuardian, error)
	ListMemberGuardians(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, memberID *uuid.UUID) ([]*ent.MemberGuardian, *base.ResponsePaginate, error)
	UpdateMemberGuardianByID(ctx context.Context, id uuid.UUID, data *ent.MemberGuardianUpdate) (*ent.MemberGuardian, error)
	SoftDeleteMemberGuardianByID(ctx context.Context, id uuid.UUID) error
}

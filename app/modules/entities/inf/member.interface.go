package entitiesinf

import (
	"context"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type MemberEntity interface {
	CreateMember(ctx context.Context, schoolID uuid.UUID, email string, password string, role ent.MemberRole, isActive bool, lastLogin *time.Time) (*ent.Member, error)
	GetMemberByEmail(ctx context.Context, email string) (*ent.Member, error)
	GetMemberByID(ctx context.Context, id uuid.UUID) (*ent.Member, error)
	ListMembers(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, role *ent.MemberRole) ([]*ent.Member, *base.ResponsePaginate, error)
	UpdateMemberByID(ctx context.Context, id uuid.UUID, schoolID *uuid.UUID, email *string, password *string, role *ent.MemberRole, isActive *bool, lastLogin *time.Time) (*ent.Member, error)
	SoftDeleteMemberByID(ctx context.Context, id uuid.UUID) error
}

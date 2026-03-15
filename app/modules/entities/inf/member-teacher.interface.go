package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type MemberTeacherEntity interface {
	CreateMemberTeacher(ctx context.Context, teacher *ent.MemberTeacher) (*ent.MemberTeacher, error)
	RegisterTeacher(ctx context.Context, data *ent.TeacherRegistrationInput) (*ent.TeacherRegistrationResult, error)
	GetMemberTeacherByID(ctx context.Context, id uuid.UUID) (*ent.MemberTeacher, error)
	ListMemberTeachers(ctx context.Context, req *base.RequestPaginate, isActive *bool, memberID *uuid.UUID, departmentID *uuid.UUID) ([]*ent.MemberTeacher, *base.ResponsePaginate, error)
	UpdateMemberTeacherByID(ctx context.Context, id uuid.UUID, req *ent.MemberTeacherUpdate) (*ent.MemberTeacher, error)
	SoftDeleteMemberTeacherByID(ctx context.Context, id uuid.UUID) error
}

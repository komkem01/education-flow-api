package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type MemberStudentEntity interface {
	CreateMemberStudent(ctx context.Context, data *ent.MemberStudent) (*ent.MemberStudent, error)
	GetMemberStudentByID(ctx context.Context, id uuid.UUID) (*ent.MemberStudent, error)
	ListMemberStudents(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, advisorTeacherID *uuid.UUID) ([]*ent.MemberStudent, *base.ResponsePaginate, error)
	UpdateMemberStudentByID(ctx context.Context, id uuid.UUID, data *ent.MemberStudentUpdate) (*ent.MemberStudent, error)
	SoftDeleteMemberStudentByID(ctx context.Context, id uuid.UUID) error
}

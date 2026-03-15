package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type StudentGuardianEntity interface {
	CreateStudentGuardian(ctx context.Context, data *ent.StudentGuardian) (*ent.StudentGuardian, error)
	GetStudentGuardianByID(ctx context.Context, id uuid.UUID) (*ent.StudentGuardian, error)
	ListStudentGuardians(ctx context.Context, req *base.RequestPaginate, studentID *uuid.UUID, guardianID *uuid.UUID, isMainGuardian *bool) ([]*ent.StudentGuardian, *base.ResponsePaginate, error)
	UpdateStudentGuardianByID(ctx context.Context, id uuid.UUID, data *ent.StudentGuardianUpdate) (*ent.StudentGuardian, error)
	SoftDeleteStudentGuardianByID(ctx context.Context, id uuid.UUID) error
}

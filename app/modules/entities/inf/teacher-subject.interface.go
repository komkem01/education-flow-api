package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type TeacherSubjectEntity interface {
	CreateTeacherSubject(ctx context.Context, data *ent.TeacherSubject) (*ent.TeacherSubject, error)
	GetTeacherSubjectByID(ctx context.Context, id uuid.UUID) (*ent.TeacherSubject, error)
	ListTeacherSubjects(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, subjectID *uuid.UUID, role *ent.TeacherSubjectRole, isActive *bool) ([]*ent.TeacherSubject, *base.ResponsePaginate, error)
	UpdateTeacherSubjectByID(ctx context.Context, id uuid.UUID, data *ent.TeacherSubjectUpdate) (*ent.TeacherSubject, error)
	SoftDeleteTeacherSubjectByID(ctx context.Context, id uuid.UUID) error
}

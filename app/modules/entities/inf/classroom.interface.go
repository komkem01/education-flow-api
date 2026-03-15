package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type ClassroomEntity interface {
	CreateClassroom(ctx context.Context, data *ent.Classroom) (*ent.Classroom, error)
	GetClassroomByID(ctx context.Context, id uuid.UUID) (*ent.Classroom, error)
	ListClassrooms(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, academicYearID *uuid.UUID, homeroomTeacherID *uuid.UUID) ([]*ent.Classroom, *base.ResponsePaginate, error)
	UpdateClassroomByID(ctx context.Context, id uuid.UUID, data *ent.ClassroomUpdate) (*ent.Classroom, error)
	SoftDeleteClassroomByID(ctx context.Context, id uuid.UUID) error
}

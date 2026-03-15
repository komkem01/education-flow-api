package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type StudentEnrollmentEntity interface {
	CreateStudentEnrollment(ctx context.Context, data *ent.StudentEnrollment) (*ent.StudentEnrollment, error)
	GetStudentEnrollmentByID(ctx context.Context, id uuid.UUID) (*ent.StudentEnrollment, error)
	ListStudentEnrollments(ctx context.Context, req *base.RequestPaginate, studentID *uuid.UUID, schoolID *uuid.UUID, academicYearID *uuid.UUID, classroomID *uuid.UUID, status *ent.StudentEnrollmentStatus, enrollmentType *ent.EnrollmentType) ([]*ent.StudentEnrollment, *base.ResponsePaginate, error)
	UpdateStudentEnrollmentByID(ctx context.Context, id uuid.UUID, data *ent.StudentEnrollmentUpdate) (*ent.StudentEnrollment, error)
	SoftDeleteStudentEnrollmentByID(ctx context.Context, id uuid.UUID) error
}

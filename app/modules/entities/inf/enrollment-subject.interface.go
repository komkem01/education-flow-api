package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type EnrollmentSubjectEntity interface {
	CreateEnrollmentSubject(ctx context.Context, data *ent.EnrollmentSubject) (*ent.EnrollmentSubject, error)
	GetEnrollmentSubjectByID(ctx context.Context, id uuid.UUID) (*ent.EnrollmentSubject, error)
	ListEnrollmentSubjects(ctx context.Context, req *base.RequestPaginate, enrollmentID *uuid.UUID, subjectID *uuid.UUID, teacherID *uuid.UUID, status *ent.StudentEnrollmentStatus, isPrimary *bool) ([]*ent.EnrollmentSubject, *base.ResponsePaginate, error)
	UpdateEnrollmentSubjectByID(ctx context.Context, id uuid.UUID, data *ent.EnrollmentSubjectUpdate) (*ent.EnrollmentSubject, error)
	SoftDeleteEnrollmentSubjectByID(ctx context.Context, id uuid.UUID) error
}

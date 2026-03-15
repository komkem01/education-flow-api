package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type TeacherLicenseEntity interface {
	CreateTeacherLicense(ctx context.Context, data *ent.TeacherLicense) (*ent.TeacherLicense, error)
	GetTeacherLicenseByID(ctx context.Context, id uuid.UUID) (*ent.TeacherLicense, error)
	ListTeacherLicenses(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, status *ent.TeacherLicenseStatus) ([]*ent.TeacherLicense, *base.ResponsePaginate, error)
	UpdateTeacherLicenseByID(ctx context.Context, id uuid.UUID, data *ent.TeacherLicenseUpdate) (*ent.TeacherLicense, error)
	SoftDeleteTeacherLicenseByID(ctx context.Context, id uuid.UUID) error
}

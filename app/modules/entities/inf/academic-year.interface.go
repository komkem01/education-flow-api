package entitiesinf

import (
	"context"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type AcademicYearEntity interface {
	CreateAcademicYear(ctx context.Context, schoolID uuid.UUID, year string, startDate time.Time, endDate time.Time, isActive bool) (*ent.AcademicYear, error)
	GetAcademicYearByID(ctx context.Context, id uuid.UUID) (*ent.AcademicYear, error)
	ListAcademicYears(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID) ([]*ent.AcademicYear, *base.ResponsePaginate, error)
	UpdateAcademicYearByID(ctx context.Context, id uuid.UUID, schoolID *uuid.UUID, year *string, startDate *time.Time, endDate *time.Time, isActive *bool) (*ent.AcademicYear, error)
	SoftDeleteAcademicYearByID(ctx context.Context, id uuid.UUID) error
}

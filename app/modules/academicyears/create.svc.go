package academicyears

import (
	"context"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, schoolID uuid.UUID, year string, startDate time.Time, endDate time.Time, isActive bool) (*ent.AcademicYear, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "academic_years.service.create")
	defer span.End()

	item, err := s.db.CreateAcademicYear(ctx, schoolID, year, startDate, endDate, isActive)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

func (s *Service) CreateAcademicYearService(ctx context.Context, schoolID uuid.UUID, year string, startDate time.Time, endDate time.Time, isActive bool) (*ent.AcademicYear, error) {
	return s.Create(ctx, schoolID, year, startDate, endDate, isActive)
}

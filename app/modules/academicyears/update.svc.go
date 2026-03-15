package academicyears

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, schoolID *uuid.UUID, year *string, startDate *time.Time, endDate *time.Time, isActive *bool) (*ent.AcademicYear, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "academic_years.service.update")
	defer span.End()

	if schoolID == nil && year == nil && startDate == nil && endDate == nil && isActive == nil {
		return nil, fmt.Errorf("%w", ErrAcademicYearConditionFail)
	}

	item, err := s.db.UpdateAcademicYearByID(ctx, id, schoolID, year, startDate, endDate, isActive)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, schoolID *uuid.UUID, year *string, startDate *time.Time, endDate *time.Time, isActive *bool) (*ent.AcademicYear, error) {
	return s.Update(ctx, id, schoolID, year, startDate, endDate, isActive)
}

package academicyears

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.AcademicYear, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "academic_years.service.get_by_id")
	defer span.End()

	item, err := s.db.GetAcademicYearByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*ent.AcademicYear, error) {
	return s.GetByID(ctx, id)
}

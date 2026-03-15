package academicyears

import (
	"context"

	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "academic_years.service.delete")
	defer span.End()

	if err := s.db.SoftDeleteAcademicYearByID(ctx, id); err != nil {
		return normalizeServiceError(err)
	}

	return nil
}

func (s *Service) DeleteService(ctx context.Context, id uuid.UUID) error {
	return s.Delete(ctx, id)
}

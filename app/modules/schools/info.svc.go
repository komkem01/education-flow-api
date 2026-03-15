package schools

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.School, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "schools.service.get_by_id")
	defer span.End()

	school, err := s.db.GetSchoolByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return school, nil
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*ent.School, error) {
	return s.GetByID(ctx, id)
}

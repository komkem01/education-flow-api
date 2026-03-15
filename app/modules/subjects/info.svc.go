package subjects

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.Subject, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "subjects.service.get_by_id")
	defer span.End()

	item, err := s.db.GetSubjectByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

package members

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.Member, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "members.service.get_by_id")
	defer span.End()

	item, err := s.db.GetMemberByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*ent.Member, error) {
	return s.GetByID(ctx, id)
}

package prefixes

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, genderID uuid.UUID, name string, isActive bool) (*ent.Prefix, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "prefixes.service.create")
	defer span.End()

	prefix, err := s.db.CreatePrefix(ctx, genderID, name, isActive)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return prefix, nil
}

func (s *Service) CreatePrefixService(ctx context.Context, genderID uuid.UUID, name string, isActive bool) (*ent.Prefix, error) {
	return s.Create(ctx, genderID, name, isActive)
}

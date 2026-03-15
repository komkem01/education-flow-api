package genders

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
)

func (s *Service) Create(ctx context.Context, name string, isActive bool) (*ent.Gender, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "genders.service.create")
	defer span.End()

	gender, err := s.db.CreateGender(ctx, name, isActive)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return gender, nil
}

func (s *Service) CreateGenderService(ctx context.Context, name string, isActive bool) (*ent.Gender, error) {
	return s.Create(ctx, name, isActive)
}

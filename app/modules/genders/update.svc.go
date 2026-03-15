package genders

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, name *string, isActive *bool) (*ent.Gender, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "genders.service.update")
	defer span.End()

	if name == nil && isActive == nil {
		return nil, fmt.Errorf("%w", ErrGenderConditionFail)
	}

	gender, err := s.db.UpdateGenderByID(ctx, id, name, isActive)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return gender, nil
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, name *string, isActive *bool) (*ent.Gender, error) {
	return s.Update(ctx, id, name, isActive)
}

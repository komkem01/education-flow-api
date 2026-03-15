package prefixes

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, genderID *uuid.UUID, name *string, isActive *bool) (*ent.Prefix, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "prefixes.service.update")
	defer span.End()

	if genderID == nil && name == nil && isActive == nil {
		return nil, fmt.Errorf("%w", ErrPrefixConditionFail)
	}

	prefix, err := s.db.UpdatePrefixByID(ctx, id, genderID, name, isActive)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return prefix, nil
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, genderID *uuid.UUID, name *string, isActive *bool) (*ent.Prefix, error) {
	return s.Update(ctx, id, genderID, name, isActive)
}

package schools

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, name *string, logoURL *string, themeColor *string, address *string, description *string) (*ent.School, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "schools.service.update")
	defer span.End()

	if name == nil && logoURL == nil && themeColor == nil && address == nil && description == nil {
		return nil, fmt.Errorf("%w", ErrSchoolConditionFail)
	}

	school, err := s.db.UpdateSchoolByID(ctx, id, name, logoURL, themeColor, address, description)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return school, nil
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, name *string, logoURL *string, themeColor *string, address *string, description *string) (*ent.School, error) {
	return s.Update(ctx, id, name, logoURL, themeColor, address, description)
}

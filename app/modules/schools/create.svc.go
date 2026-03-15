package schools

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
)

func (s *Service) Create(ctx context.Context, name string, logoURL *string, themeColor *string, address *string, description *string) (*ent.School, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "schools.service.create")
	defer span.End()

	school, err := s.db.CreateSchool(ctx, name, logoURL, themeColor, address, description)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return school, nil
}

func (s *Service) CreateSchoolService(ctx context.Context, name string, logoURL *string, themeColor *string, address *string, description *string) (*ent.School, error) {
	return s.Create(ctx, name, logoURL, themeColor, address, description)
}

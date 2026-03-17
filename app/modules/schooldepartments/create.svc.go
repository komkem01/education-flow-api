package schooldepartments

import (
	"context"
	"fmt"
	"strings"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
)

func (s *Service) Create(ctx context.Context, data *ent.SchoolDepartment) (*ent.SchoolDepartment, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "schooldepartments.service.create")
	defer span.End()

	if data == nil {
		return nil, fmt.Errorf("%w", ErrSchoolDepartmentConditionFail)
	}

	if code := strings.ToUpper(strings.TrimSpace(data.Code)); code != "" {
		data.Code = code
		item, err := s.db.CreateSchoolDepartment(ctx, data)
		if err != nil {
			return nil, normalizeServiceError(err)
		}
		return item, nil
	}

	const maxCodeRetry = 10
	for i := 0; i < maxCodeRetry; i++ {
		generatedCode, err := utils.GenerateNumericCode("SD", 6)
		if err != nil {
			return nil, err
		}
		data.Code = generatedCode

		item, err := s.db.CreateSchoolDepartment(ctx, data)
		if err == nil {
			return item, nil
		}
		if isDuplicateKeyError(err) {
			continue
		}
		return nil, normalizeServiceError(err)
	}

	return nil, fmt.Errorf("%w", ErrSchoolDepartmentDuplicate)
}

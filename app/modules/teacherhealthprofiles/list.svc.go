package teacherhealthprofiles

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, memberTeacherID *uuid.UUID, bloodType *string) ([]*ent.TeacherHealthProfile, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacherhealthprofiles.service.list")
	defer span.End()

	var bloodTypeVal *string
	if bloodType != nil {
		parsed, ok := parseBloodType(*bloodType)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrTeacherHealthProfileConditionFail)
		}
		bloodTypeVal = &parsed
	}

	items, page, err := s.db.ListTeacherHealthProfiles(ctx, req, memberTeacherID, bloodTypeVal)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrTeacherHealthProfileConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

package studenthealthprofiles

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, studentID *uuid.UUID, bloodType *string) ([]*ent.StudentHealthProfile, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studenthealthprofiles.service.list")
	defer span.End()

	items, page, err := s.db.ListStudentHealthProfiles(ctx, req, studentID, bloodType)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrStudentHealthProfileConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

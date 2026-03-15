package studentguardians

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, studentID *uuid.UUID, guardianID *uuid.UUID, isMainGuardian *bool) ([]*ent.StudentGuardian, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentguardians.service.list")
	defer span.End()

	items, page, err := s.db.ListStudentGuardians(ctx, req, studentID, guardianID, isMainGuardian)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrStudentGuardianConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

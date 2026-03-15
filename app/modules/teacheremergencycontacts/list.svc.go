package teacheremergencycontacts

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, memberTeacherID *uuid.UUID, isPrimary *bool) ([]*ent.TeacherEmergencyContact, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacheremergencycontacts.service.list")
	defer span.End()

	items, page, err := s.db.ListTeacherEmergencyContacts(ctx, req, memberTeacherID, isPrimary)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrTeacherEmergencyContactConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

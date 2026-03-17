package memberteachers

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) ListAddresses(ctx context.Context, memberTeacherID uuid.UUID) ([]*ent.TeacherAddress, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberteachers.service.list_addresses")
	defer span.End()

	if _, err := s.db.GetMemberTeacherByID(ctx, memberTeacherID); err != nil {
		return nil, normalizeServiceError(err)
	}

	items, err := s.db.ListTeacherAddressesByMemberTeacherID(ctx, memberTeacherID)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return items, nil
}

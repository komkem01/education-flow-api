package schools

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "schools.service.delete")
	defer span.End()

	role := ent.MemberRoleAdmin
	pageReq := &base.RequestPaginate{Page: 1, Size: 100}
	admins, _, err := s.member.ListMembers(ctx, pageReq, nil, &id, &role)
	if err != nil {
		return err
	}

	for _, admin := range admins {
		if admin == nil {
			continue
		}
		if err := s.member.SoftDeleteMemberByID(ctx, admin.ID); err != nil {
			return err
		}
	}

	if err := s.db.SoftDeleteSchoolByID(ctx, id); err != nil {
		return normalizeServiceError(err)
	}

	return nil
}

func (s *Service) DeleteService(ctx context.Context, id uuid.UUID) error {
	return s.Delete(ctx, id)
}

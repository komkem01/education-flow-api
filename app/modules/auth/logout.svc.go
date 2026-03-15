package auth

import (
	"context"

	"eduflow/app/utils"
)

func (s *Service) Logout(ctx context.Context, bearer string) error {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "auth.service.logout")
	defer span.End()

	user, err := s.resolveCurrentUser(ctx, bearer)
	if err != nil {
		return err
	}

	return normalizeServiceError(s.session.DeleteAuthSessionByToken(ctx, user.Claims.JTI))
}

func (s *Service) LogoutAll(ctx context.Context, bearer string) error {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "auth.service.logout_all")
	defer span.End()

	user, err := s.resolveCurrentUser(ctx, bearer)
	if err != nil {
		return err
	}

	return normalizeServiceError(s.session.DeleteAuthSessionsByMemberID(ctx, user.Member.ID))
}

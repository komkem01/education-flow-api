package auth

import (
	"context"
	"time"

	"eduflow/app/utils"
)

func (s *Service) Me(ctx context.Context, bearer string) (*MeResponse, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "auth.service.me")
	defer span.End()

	user, err := s.resolveCurrentUser(ctx, bearer)
	if err != nil {
		return nil, err
	}

	return &MeResponse{
		ID:        user.Member.ID.String(),
		SchoolID:  user.Member.SchoolID.String(),
		Email:     user.Member.Email,
		Role:      string(user.Member.Role),
		IsActive:  user.Member.IsActive,
		LastLogin: user.Member.LastLogin,
		ExpireAt:  time.Unix(user.Claims.EXP, 0),
	}, nil
}

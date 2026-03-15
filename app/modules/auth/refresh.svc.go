package auth

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "auth.service.refresh")
	defer span.End()

	if refreshToken == "" {
		return nil, fmt.Errorf("%w", ErrAuthConditionFail)
	}

	hash := hashToken(refreshToken)
	session, err := s.session.GetAuthSessionByRefreshToken(ctx, hash)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	now := time.Now()
	if session.RefreshExpireAt.Before(now) {
		return nil, fmt.Errorf("%w", ErrAuthUnauthorized)
	}

	member, err := s.member.GetMemberByID(ctx, session.MemberID)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	if !member.IsActive {
		return nil, fmt.Errorf("%w", ErrAuthUnauthorized)
	}

	jti := uuid.NewString()
	expireAt := now.Add(time.Duration(s.conf.AccessTokenTTLSeconds) * time.Second)
	refreshExpireAt := now.Add(time.Duration(s.conf.RefreshTokenTTLSeconds) * time.Second)
	token, err := s.createAccessToken(&accessClaims{
		Sub:      member.ID.String(),
		Role:     string(member.Role),
		SchoolID: member.SchoolID.String(),
		JTI:      jti,
		IAT:      now.Unix(),
		EXP:      expireAt.Unix(),
	})
	if err != nil {
		return nil, err
	}

	newRefresh, err := newRandomToken()
	if err != nil {
		return nil, err
	}

	_, err = s.session.UpdateAuthSessionByID(ctx, session.ID, &ent.AuthSessionUpdate{
		Token:           &jti,
		RefreshToken:    ptrString(hashToken(newRefresh)),
		ExpireAt:        &expireAt,
		RefreshExpireAt: &refreshExpireAt,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return &TokenResponse{Token: token, RefreshToken: newRefresh, ExpireAt: expireAt}, nil
}

func ptrString(v string) *string {
	return &v
}

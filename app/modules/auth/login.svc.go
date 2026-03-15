package auth

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/hashing"

	"github.com/google/uuid"
)

func (s *Service) Login(ctx context.Context, email string, password string) (*TokenResponse, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "auth.service.login")
	defer span.End()

	if email == "" || password == "" {
		return nil, fmt.Errorf("%w", ErrAuthConditionFail)
	}

	member, err := s.member.GetMemberByEmail(ctx, email)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	if !member.IsActive {
		return nil, fmt.Errorf("%w", ErrAuthUnauthorized)
	}
	if !hashing.CheckPasswordHash([]byte(member.Password), []byte(password)) {
		return nil, fmt.Errorf("%w", ErrAuthUnauthorized)
	}

	now := time.Now()
	expireAt := now.Add(time.Duration(s.conf.AccessTokenTTLSeconds) * time.Second)
	refreshExpireAt := now.Add(time.Duration(s.conf.RefreshTokenTTLSeconds) * time.Second)
	jti := uuid.NewString()

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

	refreshToken, err := newRandomToken()
	if err != nil {
		return nil, err
	}

	_, err = s.session.CreateAuthSession(ctx, &ent.AuthSession{
		MemberID:        member.ID,
		Token:           jti,
		RefreshToken:    hashToken(refreshToken),
		ExpireAt:        expireAt,
		RefreshExpireAt: refreshExpireAt,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	_, _ = s.member.UpdateMemberByID(ctx, member.ID, nil, nil, nil, nil, nil, &now)

	return &TokenResponse{Token: token, RefreshToken: refreshToken, ExpireAt: expireAt}, nil
}

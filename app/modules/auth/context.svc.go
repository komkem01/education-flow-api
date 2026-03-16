package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"

	"github.com/google/uuid"
)

type CurrentUser struct {
	Member *ent.Member
	Claims *accessClaims
}

func (s *Service) resolveCurrentUser(ctx context.Context, bearer string) (*CurrentUser, error) {
	token := extractBearerToken(bearer)
	if token == "" {
		token = strings.TrimSpace(bearer)
	}
	if token == "" {
		return nil, fmt.Errorf("%w", ErrAuthUnauthorized)
	}

	claims, err := s.parseAccessToken(token)
	if err != nil {
		return nil, err
	}

	session, err := s.session.GetAuthSessionByToken(ctx, claims.JTI)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	if session.ExpireAt.Before(time.Now()) {
		return nil, fmt.Errorf("%w", ErrAuthUnauthorized)
	}

	memberID, err := uuid.Parse(claims.Sub)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAuthUnauthorized)
	}
	if session.MemberID != memberID {
		return nil, fmt.Errorf("%w", ErrAuthUnauthorized)
	}

	member, err := s.member.GetMemberByID(ctx, memberID)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	if !member.IsActive {
		return nil, fmt.Errorf("%w", ErrAuthUnauthorized)
	}

	return &CurrentUser{Member: member, Claims: claims}, nil
}

func extractBearerToken(v string) string {
	if v == "" {
		return ""
	}
	parts := strings.SplitN(strings.TrimSpace(v), " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

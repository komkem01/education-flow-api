package entitiesinf

import (
	"context"

	"eduflow/app/modules/entities/ent"

	"github.com/google/uuid"
)

type AuthSessionEntity interface {
	CreateAuthSession(ctx context.Context, data *ent.AuthSession) (*ent.AuthSession, error)
	GetAuthSessionByToken(ctx context.Context, token string) (*ent.AuthSession, error)
	GetAuthSessionByRefreshToken(ctx context.Context, refreshToken string) (*ent.AuthSession, error)
	UpdateAuthSessionByID(ctx context.Context, id uuid.UUID, data *ent.AuthSessionUpdate) (*ent.AuthSession, error)
	DeleteAuthSessionByToken(ctx context.Context, token string) error
	DeleteAuthSessionsByMemberID(ctx context.Context, memberID uuid.UUID) error
}

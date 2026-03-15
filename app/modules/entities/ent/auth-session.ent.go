package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type AuthSession struct {
	bun.BaseModel `bun:"table:auth_sessions,alias:as"`

	ID              uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberID        uuid.UUID `bun:"member_id,notnull,type:uuid"`
	Token           string    `bun:"token,notnull"`
	RefreshToken    string    `bun:"refresh_token,notnull"`
	ExpireAt        time.Time `bun:"expire_at,notnull"`
	RefreshExpireAt time.Time `bun:"refresh_expire_at,notnull"`
	CreatedAt       time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt       time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type AuthSessionUpdate struct {
	Token           *string
	RefreshToken    *string
	ExpireAt        *time.Time
	RefreshExpireAt *time.Time
}

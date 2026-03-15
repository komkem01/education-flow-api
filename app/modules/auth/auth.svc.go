package auth

import (
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/internal/config"
	"os"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	AccessTokenTTLSeconds  int    `conf:"default=900"`
	RefreshTokenTTLSeconds int    `conf:"default=604800"`
	JWTSecret              string `conf:"default=change-me"`
}

type Options struct {
	*config.Config[Config]
	tracer  trace.Tracer
	member  entitiesinf.MemberEntity
	session entitiesinf.AuthSessionEntity
}

type Service struct {
	tracer  trace.Tracer
	member  entitiesinf.MemberEntity
	session entitiesinf.AuthSessionEntity
	conf    Config
}

type TokenResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpireAt     time.Time `json:"expire_at"`
}

type MeResponse struct {
	ID        string     `json:"id"`
	SchoolID  string     `json:"school_id"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	IsActive  bool       `json:"is_active"`
	LastLogin *time.Time `json:"last_login"`
	ExpireAt  time.Time  `json:"expire_at"`
}

func newService(opt *Options) *Service {
	conf := Config{}
	if opt != nil && opt.Config != nil {
		conf = *opt.Config.Val
	}
	if conf.AccessTokenTTLSeconds <= 0 {
		conf.AccessTokenTTLSeconds = 900
	}
	if conf.RefreshTokenTTLSeconds <= 0 {
		conf.RefreshTokenTTLSeconds = 604800
	}
	if conf.JWTSecret == "" {
		conf.JWTSecret = "change-me"
	}
	if envSecret := os.Getenv("JWT_SECRET"); envSecret != "" {
		conf.JWTSecret = envSecret
	}

	return &Service{tracer: opt.tracer, member: opt.member, session: opt.session, conf: conf}
}

package members

import (
	"context"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/hashing"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, schoolID uuid.UUID, email string, password string, role ent.MemberRole, isActive bool, lastLogin *time.Time) (*ent.Member, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "members.service.create")
	defer span.End()

	if strings.TrimSpace(password) == "" {
		return nil, ErrMemberConditionFail
	}
	if role == "" {
		role = ent.MemberRoleAdmin
	}

	hashed, err := hashing.HashPassword(password)
	if err != nil {
		return nil, err
	}

	item, err := s.db.CreateMember(ctx, schoolID, email, string(hashed), role, isActive, lastLogin)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

func (s *Service) CreateMemberService(ctx context.Context, schoolID uuid.UUID, email string, password string, role ent.MemberRole, isActive bool, lastLogin *time.Time) (*ent.Member, error) {
	return s.Create(ctx, schoolID, email, password, role, isActive, lastLogin)
}

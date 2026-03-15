package members

import (
	"context"
	"fmt"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/hashing"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, schoolID *uuid.UUID, email *string, password *string, role *ent.MemberRole, isActive *bool, lastLogin *time.Time) (*ent.Member, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "members.service.update")
	defer span.End()

	if schoolID == nil && email == nil && password == nil && role == nil && isActive == nil && lastLogin == nil {
		return nil, fmt.Errorf("%w", ErrMemberConditionFail)
	}

	var hashedPassword *string
	if password != nil {
		if strings.TrimSpace(*password) == "" {
			return nil, fmt.Errorf("%w", ErrMemberConditionFail)
		}
		hashed, err := hashing.HashPassword(*password)
		if err != nil {
			return nil, err
		}
		hashedStr := string(hashed)
		hashedPassword = &hashedStr
	}

	item, err := s.db.UpdateMemberByID(ctx, id, schoolID, email, hashedPassword, role, isActive, lastLogin)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, schoolID *uuid.UUID, email *string, password *string, role *ent.MemberRole, isActive *bool, lastLogin *time.Time) (*ent.Member, error) {
	return s.Update(ctx, id, schoolID, email, password, role, isActive, lastLogin)
}

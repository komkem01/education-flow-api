package memberguardians

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, memberID *uuid.UUID, schoolID uuid.UUID, genderID uuid.UUID, prefixID uuid.UUID, firstNameTH string, lastNameTH string, firstNameEN *string, lastNameEN *string, citizenID *string, phone *string, isActive bool) (*ent.MemberGuardian, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberguardians.service.create")
	defer span.End()

	item, err := s.db.CreateMemberGuardian(ctx, &ent.MemberGuardian{
		MemberID:    memberID,
		SchoolID:    schoolID,
		GenderID:    genderID,
		PrefixID:    prefixID,
		FirstNameTH: firstNameTH,
		LastNameTH:  lastNameTH,
		FirstNameEN: firstNameEN,
		LastNameEN:  lastNameEN,
		CitizenID:   citizenID,
		Phone:       phone,
		IsActive:    isActive,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

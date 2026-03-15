package memberguardians

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.MemberGuardian, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberguardians.service.update")
	defer span.End()

	if req.MemberID == nil && req.SchoolID == nil && req.GenderID == nil && req.PrefixID == nil && req.FirstNameTH == nil && req.LastNameTH == nil && req.FirstNameEN == nil && req.LastNameEN == nil && req.CitizenID == nil && req.Phone == nil && req.IsActive == nil {
		return nil, fmt.Errorf("%w", ErrMemberGuardianConditionFail)
	}

	payload := &ent.MemberGuardianUpdate{
		FirstNameTH: req.FirstNameTH,
		LastNameTH:  req.LastNameTH,
		FirstNameEN: req.FirstNameEN,
		LastNameEN:  req.LastNameEN,
		CitizenID:   req.CitizenID,
		Phone:       req.Phone,
		IsActive:    req.IsActive,
	}

	if req.MemberID != nil {
		parsed, err := uuid.Parse(*req.MemberID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberGuardianConditionFail)
		}
		payload.MemberID = &parsed
	}
	if req.SchoolID != nil {
		parsed, err := uuid.Parse(*req.SchoolID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberGuardianConditionFail)
		}
		payload.SchoolID = &parsed
	}
	if req.GenderID != nil {
		parsed, err := uuid.Parse(*req.GenderID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberGuardianConditionFail)
		}
		payload.GenderID = &parsed
	}
	if req.PrefixID != nil {
		parsed, err := uuid.Parse(*req.PrefixID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrMemberGuardianConditionFail)
		}
		payload.PrefixID = &parsed
	}

	item, err := s.db.UpdateMemberGuardianByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

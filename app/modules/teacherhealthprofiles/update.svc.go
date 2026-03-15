package teacherhealthprofiles

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.TeacherHealthProfile, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacherhealthprofiles.service.update")
	defer span.End()

	if req.MemberTeacherID == nil && req.BloodType == nil && req.AllergyInfo == nil && req.ChronicDisease == nil && req.MedicationNote == nil && req.FitnessForWorkNote == nil {
		return nil, fmt.Errorf("%w", ErrTeacherHealthProfileConditionFail)
	}

	payload := &ent.TeacherHealthProfileUpdate{
		AllergyInfo:        req.AllergyInfo,
		ChronicDisease:     req.ChronicDisease,
		MedicationNote:     req.MedicationNote,
		FitnessForWorkNote: req.FitnessForWorkNote,
	}

	if req.MemberTeacherID != nil {
		parsed, err := uuid.Parse(*req.MemberTeacherID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrTeacherHealthProfileConditionFail)
		}
		payload.MemberTeacherID = &parsed
	}
	if req.BloodType != nil {
		parsed, ok := parseBloodType(*req.BloodType)
		if !ok {
			return nil, fmt.Errorf("%w", ErrTeacherHealthProfileConditionFail)
		}
		payload.BloodType = &parsed
	}

	item, err := s.db.UpdateTeacherHealthProfileByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

package teacherhealthprofiles

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, memberTeacherID uuid.UUID, bloodType *string, allergyInfo *string, chronicDisease *string, medicationNote *string, fitnessForWorkNote *string) (*ent.TeacherHealthProfile, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacherhealthprofiles.service.create")
	defer span.End()

	var bloodTypeVal *string
	if bloodType != nil {
		parsed, ok := parseBloodType(*bloodType)
		if !ok {
			return nil, fmt.Errorf("%w", ErrTeacherHealthProfileConditionFail)
		}
		bloodTypeVal = &parsed
	}

	item, err := s.db.CreateTeacherHealthProfile(ctx, &ent.TeacherHealthProfile{
		MemberTeacherID:    memberTeacherID,
		BloodType:          bloodTypeVal,
		AllergyInfo:        allergyInfo,
		ChronicDisease:     chronicDisease,
		MedicationNote:     medicationNote,
		FitnessForWorkNote: fitnessForWorkNote,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

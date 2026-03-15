package studenthealthprofiles

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.StudentHealthProfile, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studenthealthprofiles.service.update")
	defer span.End()

	if req.StudentID == nil && req.BloodType == nil && req.AllergyInfo == nil && req.ChronicDisease == nil && req.MedicalNote == nil {
		return nil, fmt.Errorf("%w", ErrStudentHealthProfileConditionFail)
	}

	payload := &ent.StudentHealthProfileUpdate{
		BloodType:      req.BloodType,
		AllergyInfo:    req.AllergyInfo,
		ChronicDisease: req.ChronicDisease,
		MedicalNote:    req.MedicalNote,
	}

	if req.StudentID != nil {
		parsed, err := uuid.Parse(*req.StudentID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrStudentHealthProfileConditionFail)
		}
		payload.StudentID = &parsed
	}

	item, err := s.db.UpdateStudentHealthProfileByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

package studentprofiles

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.StudentProfile, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentprofiles.service.update")
	defer span.End()

	if req.StudentID == nil && req.BirthDate == nil && req.Nationality == nil && req.Religion == nil && req.AddressCurrent == nil && req.AddressRegistered == nil && req.EmergencyContactName == nil && req.EmergencyContactPhone == nil {
		return nil, fmt.Errorf("%w", ErrStudentProfileConditionFail)
	}

	payload := &ent.StudentProfileUpdate{
		Nationality:           req.Nationality,
		Religion:              req.Religion,
		AddressCurrent:        req.AddressCurrent,
		AddressRegistered:     req.AddressRegistered,
		EmergencyContactName:  req.EmergencyContactName,
		EmergencyContactPhone: req.EmergencyContactPhone,
	}

	if req.StudentID != nil {
		parsed, err := uuid.Parse(*req.StudentID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrStudentProfileConditionFail)
		}
		payload.StudentID = &parsed
	}
	if req.BirthDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrStudentProfileConditionFail)
		}
		payload.BirthDate = &parsed
	}

	item, err := s.db.UpdateStudentProfileByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

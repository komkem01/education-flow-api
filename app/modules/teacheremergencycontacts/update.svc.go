package teacheremergencycontacts

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.TeacherEmergencyContact, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacheremergencycontacts.service.update")
	defer span.End()

	if req.EmergencyContactName == nil && req.Relationship == nil && req.PhonePrimary == nil && req.PhoneSecondary == nil && req.CanDecideMedical == nil && req.IsPrimary == nil {
		return nil, fmt.Errorf("%w", ErrTeacherEmergencyContactConditionFail)
	}

	payload := &ent.TeacherEmergencyContactUpdate{
		EmergencyContactName: req.EmergencyContactName,
		Relationship:         req.Relationship,
		PhonePrimary:         req.PhonePrimary,
		PhoneSecondary:       req.PhoneSecondary,
		CanDecideMedical:     req.CanDecideMedical,
		IsPrimary:            req.IsPrimary,
	}

	updated, err := s.db.UpdateTeacherEmergencyContactByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return updated, nil
}

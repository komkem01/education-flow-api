package studentguardians

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.StudentGuardian, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentguardians.service.update")
	defer span.End()

	if req.StudentID == nil && req.GuardianID == nil && req.Relationship == nil && req.IsMainGuardian == nil && req.CanPickup == nil && req.IsEmergencyContact == nil && req.Note == nil {
		return nil, fmt.Errorf("%w", ErrStudentGuardianConditionFail)
	}

	payload := &ent.StudentGuardianUpdate{
		Relationship:       req.Relationship,
		IsMainGuardian:     req.IsMainGuardian,
		CanPickup:          req.CanPickup,
		IsEmergencyContact: req.IsEmergencyContact,
		Note:               req.Note,
	}

	if req.StudentID != nil {
		parsed, err := uuid.Parse(*req.StudentID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrStudentGuardianConditionFail)
		}
		payload.StudentID = &parsed
	}
	if req.GuardianID != nil {
		parsed, err := uuid.Parse(*req.GuardianID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrStudentGuardianConditionFail)
		}
		payload.GuardianID = &parsed
	}

	item, err := s.db.UpdateStudentGuardianByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

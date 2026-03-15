package studentguardians

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, studentID uuid.UUID, guardianID uuid.UUID, relationship string, isMainGuardian bool, canPickup bool, isEmergencyContact bool, note *string) (*ent.StudentGuardian, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentguardians.service.create")
	defer span.End()

	item, err := s.db.CreateStudentGuardian(ctx, &ent.StudentGuardian{
		StudentID:          studentID,
		GuardianID:         guardianID,
		Relationship:       relationship,
		IsMainGuardian:     isMainGuardian,
		CanPickup:          canPickup,
		IsEmergencyContact: isEmergencyContact,
		Note:               note,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

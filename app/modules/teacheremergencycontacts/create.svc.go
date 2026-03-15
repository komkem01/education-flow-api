package teacheremergencycontacts

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, memberTeacherID uuid.UUID, emergencyContactName string, relationship string, phonePrimary string, phoneSecondary *string, canDecideMedical bool, isPrimary bool) (*ent.TeacherEmergencyContact, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacheremergencycontacts.service.create")
	defer span.End()

	item, err := s.db.CreateTeacherEmergencyContact(ctx, &ent.TeacherEmergencyContact{
		MemberTeacherID:      memberTeacherID,
		EmergencyContactName: emergencyContactName,
		Relationship:         relationship,
		PhonePrimary:         phonePrimary,
		PhoneSecondary:       phoneSecondary,
		CanDecideMedical:     canDecideMedical,
		IsPrimary:            isPrimary,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

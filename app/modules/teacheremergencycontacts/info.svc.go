package teacheremergencycontacts

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.TeacherEmergencyContact, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacheremergencycontacts.service.info")
	defer span.End()

	item, err := s.db.GetTeacherEmergencyContactByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

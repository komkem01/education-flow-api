package studentprofiles

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.StudentProfile, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentprofiles.service.get_by_id")
	defer span.End()

	item, err := s.db.GetStudentProfileByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

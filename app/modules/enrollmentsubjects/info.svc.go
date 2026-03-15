package enrollmentsubjects

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.EnrollmentSubject, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "enrollmentsubjects.service.info")
	defer span.End()

	item, err := s.db.GetEnrollmentSubjectByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

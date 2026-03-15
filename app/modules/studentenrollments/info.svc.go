package studentenrollments

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.StudentEnrollment, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentenrollments.service.info")
	defer span.End()

	item, err := s.db.GetStudentEnrollmentByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

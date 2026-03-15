package teacherexperiences

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, isCurrent *bool, isActive *bool) ([]*ent.TeacherExperience, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacherexperiences.service.list")
	defer span.End()

	items, page, err := s.db.ListTeacherExperiences(ctx, req, teacherID, isCurrent, isActive)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrTeacherExperienceConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

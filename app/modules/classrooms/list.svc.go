package classrooms

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, isActive *bool, schoolID *uuid.UUID, academicYearID *uuid.UUID, homeroomTeacherID *uuid.UUID) ([]*ent.Classroom, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "classrooms.service.list")
	defer span.End()

	items, page, err := s.db.ListClassrooms(ctx, req, isActive, schoolID, academicYearID, homeroomTeacherID)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrClassroomConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

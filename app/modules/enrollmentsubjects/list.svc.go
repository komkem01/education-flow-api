package enrollmentsubjects

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, enrollmentID *uuid.UUID, subjectID *uuid.UUID, teacherID *uuid.UUID, status *string, isPrimary *bool) ([]*ent.EnrollmentSubject, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "enrollmentsubjects.service.list")
	defer span.End()

	var statusVal *ent.StudentEnrollmentStatus
	if status != nil {
		parsed, ok := parseStudentEnrollmentStatus(*status)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrEnrollmentSubjectConditionFail)
		}
		statusVal = &parsed
	}

	items, page, err := s.db.ListEnrollmentSubjects(ctx, req, enrollmentID, subjectID, teacherID, statusVal, isPrimary)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrEnrollmentSubjectConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

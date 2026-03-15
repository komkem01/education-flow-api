package studentenrollments

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, studentID *uuid.UUID, schoolID *uuid.UUID, academicYearID *uuid.UUID, classroomID *uuid.UUID, status *string, enrollmentType *string) ([]*ent.StudentEnrollment, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentenrollments.service.list")
	defer span.End()

	var statusVal *ent.StudentEnrollmentStatus
	if status != nil {
		parsed, ok := parseStudentEnrollmentStatus(*status)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
		}
		statusVal = &parsed
	}

	var enrollmentTypeVal *ent.EnrollmentType
	if enrollmentType != nil {
		parsed, ok := parseEnrollmentType(*enrollmentType)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
		}
		enrollmentTypeVal = &parsed
	}

	items, page, err := s.db.ListStudentEnrollments(ctx, req, studentID, schoolID, academicYearID, classroomID, statusVal, enrollmentTypeVal)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrStudentEnrollmentConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

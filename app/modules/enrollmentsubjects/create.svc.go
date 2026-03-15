package enrollmentsubjects

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, enrollmentID uuid.UUID, subjectID uuid.UUID, teacherID *string, isPrimary bool, status *string) (*ent.EnrollmentSubject, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "enrollmentsubjects.service.create")
	defer span.End()

	statusVal := ent.StudentEnrollmentStatusActive
	if status != nil && *status != "" {
		parsed, ok := parseStudentEnrollmentStatus(*status)
		if !ok {
			return nil, fmt.Errorf("%w", ErrEnrollmentSubjectConditionFail)
		}
		statusVal = parsed
	}

	teacherIDVal, err := parseOptionalUUID(teacherID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrEnrollmentSubjectConditionFail)
	}

	item, err := s.db.CreateEnrollmentSubject(ctx, &ent.EnrollmentSubject{
		EnrollmentID: enrollmentID,
		SubjectID:    subjectID,
		TeacherID:    teacherIDVal,
		IsPrimary:    isPrimary,
		Status:       statusVal,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

func parseOptionalUUID(v *string) (*uuid.UUID, error) {
	if v == nil || *v == "" {
		return nil, nil
	}
	parsed, err := uuid.Parse(*v)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

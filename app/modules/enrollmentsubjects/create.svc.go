package enrollmentsubjects

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(
	ctx context.Context,
	enrollmentID uuid.UUID,
	subjectID uuid.UUID,
	teacherID *string,
	isPrimary bool,
	status *string,
	midtermScore *float64,
	finalScore *float64,
	activityScore *float64,
) (*ent.EnrollmentSubject, error) {
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

	if !isValidScore(midtermScore) || !isValidScore(finalScore) || !isValidScore(activityScore) {
		return nil, fmt.Errorf("%w", ErrEnrollmentSubjectConditionFail)
	}

	item, err := s.db.CreateEnrollmentSubject(ctx, &ent.EnrollmentSubject{
		EnrollmentID:  enrollmentID,
		SubjectID:     subjectID,
		TeacherID:     teacherIDVal,
		IsPrimary:     isPrimary,
		Status:        statusVal,
		MidtermScore:  midtermScore,
		FinalScore:    finalScore,
		ActivityScore: activityScore,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	enrichScoreSummary(item)

	return item, nil
}

func isValidScore(score *float64) bool {
	if score == nil {
		return true
	}
	return *score >= 0 && *score <= 100
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

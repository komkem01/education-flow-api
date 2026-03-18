package enrollmentsubjects

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.EnrollmentSubject, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "enrollmentsubjects.service.update")
	defer span.End()

	if req.EnrollmentID == nil && req.SubjectID == nil && req.TeacherID == nil && req.IsPrimary == nil && req.Status == nil && req.MidtermScore == nil && req.FinalScore == nil && req.ActivityScore == nil {
		return nil, fmt.Errorf("%w", ErrEnrollmentSubjectConditionFail)
	}

	if !isValidScore(req.MidtermScore) || !isValidScore(req.FinalScore) || !isValidScore(req.ActivityScore) {
		return nil, fmt.Errorf("%w", ErrEnrollmentSubjectConditionFail)
	}

	payload := &ent.EnrollmentSubjectUpdate{
		IsPrimary:     req.IsPrimary,
		MidtermScore:  req.MidtermScore,
		FinalScore:    req.FinalScore,
		ActivityScore: req.ActivityScore,
	}

	var err error
	payload.EnrollmentID, err = parseOptionalUUID(req.EnrollmentID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrEnrollmentSubjectConditionFail)
	}
	payload.SubjectID, err = parseOptionalUUID(req.SubjectID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrEnrollmentSubjectConditionFail)
	}
	payload.TeacherID, err = parseOptionalUUID(req.TeacherID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrEnrollmentSubjectConditionFail)
	}

	if req.Status != nil {
		parsed, ok := parseStudentEnrollmentStatus(*req.Status)
		if !ok {
			return nil, fmt.Errorf("%w", ErrEnrollmentSubjectConditionFail)
		}
		payload.Status = &parsed
	}

	item, err := s.db.UpdateEnrollmentSubjectByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	enrichScoreSummary(item)

	return item, nil
}

package studentenrollments

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.StudentEnrollment, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentenrollments.service.update")
	defer span.End()

	if req.StudentID == nil && req.SchoolID == nil && req.AcademicYearID == nil && req.ClassroomID == nil && req.EnrolledAt == nil && req.ExitedAt == nil && req.Status == nil && req.EnrollmentType == nil && req.ExitReason == nil && req.ExitNote == nil && req.PreviousEnrollmentID == nil && req.RollNo == nil && req.ApprovedBy == nil && req.ApprovedAt == nil && req.ApprovalNote == nil && req.CreatedBy == nil && req.UpdatedBy == nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}

	payload := &ent.StudentEnrollmentUpdate{
		ExitNote:     req.ExitNote,
		RollNo:       req.RollNo,
		ApprovalNote: req.ApprovalNote,
	}

	var err error
	payload.StudentID, err = parseOptionalUUID(req.StudentID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}
	payload.SchoolID, err = parseOptionalUUID(req.SchoolID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}
	payload.AcademicYearID, err = parseOptionalUUID(req.AcademicYearID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}
	payload.ClassroomID, err = parseOptionalUUID(req.ClassroomID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}
	payload.PreviousEnrollmentID, err = parseOptionalUUID(req.PreviousEnrollmentID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}
	payload.ApprovedBy, err = parseOptionalUUID(req.ApprovedBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}
	payload.CreatedBy, err = parseOptionalUUID(req.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}
	payload.UpdatedBy, err = parseOptionalUUID(req.UpdatedBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}

	payload.EnrolledAt, err = parseOptionalDate(req.EnrolledAt)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}
	payload.ExitedAt, err = parseOptionalDate(req.ExitedAt)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}
	payload.ApprovedAt, err = parseOptionalDate(req.ApprovedAt)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}

	if payload.EnrolledAt != nil && payload.ExitedAt != nil && payload.ExitedAt.Before(*payload.EnrolledAt) {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}

	if req.Status != nil {
		parsed, ok := parseStudentEnrollmentStatus(*req.Status)
		if !ok {
			return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
		}
		payload.Status = &parsed
	}
	if req.EnrollmentType != nil {
		parsed, ok := parseEnrollmentType(*req.EnrollmentType)
		if !ok {
			return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
		}
		payload.EnrollmentType = &parsed
	}
	if req.ExitReason != nil {
		if *req.ExitReason == "" {
			payload.ExitReason = nil
		} else {
			parsed, ok := parseEnrollmentExitReason(*req.ExitReason)
			if !ok {
				return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
			}
			payload.ExitReason = &parsed
		}
	}

	item, err := s.db.UpdateStudentEnrollmentByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

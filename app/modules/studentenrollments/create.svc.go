package studentenrollments

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, studentID uuid.UUID, schoolID uuid.UUID, academicYearID uuid.UUID, classroomID uuid.UUID, enrolledAt *string, exitedAt *string, status *string, enrollmentType *string, exitReason *string, exitNote *string, previousEnrollmentID *string, rollNo *string, approvedBy *string, approvedAt *string, approvalNote *string, createdBy *string, updatedBy *string) (*ent.StudentEnrollment, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentenrollments.service.create")
	defer span.End()

	statusVal := ent.StudentEnrollmentStatusActive
	if status != nil && *status != "" {
		parsed, ok := parseStudentEnrollmentStatus(*status)
		if !ok {
			return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
		}
		statusVal = parsed
	}

	enrollmentTypeVal := ent.EnrollmentTypeNew
	if enrollmentType != nil && *enrollmentType != "" {
		parsed, ok := parseEnrollmentType(*enrollmentType)
		if !ok {
			return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
		}
		enrollmentTypeVal = parsed
	}

	var enrolledAtVal *time.Time
	if enrolledAt != nil && *enrolledAt != "" {
		parsed, err := time.Parse("2006-01-02", *enrolledAt)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
		}
		enrolledAtVal = &parsed
	}

	var exitedAtVal *time.Time
	if exitedAt != nil && *exitedAt != "" {
		parsed, err := time.Parse("2006-01-02", *exitedAt)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
		}
		exitedAtVal = &parsed
	}

	if enrolledAtVal != nil && exitedAtVal != nil && exitedAtVal.Before(*enrolledAtVal) {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}

	var exitReasonVal *ent.EnrollmentExitReason
	if exitReason != nil && *exitReason != "" {
		parsed, ok := parseEnrollmentExitReason(*exitReason)
		if !ok {
			return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
		}
		exitReasonVal = &parsed
	}

	previousEnrollmentIDVal, err := parseOptionalUUID(previousEnrollmentID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}
	approvedByVal, err := parseOptionalUUID(approvedBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}
	createdByVal, err := parseOptionalUUID(createdBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}
	updatedByVal, err := parseOptionalUUID(updatedBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}

	approvedAtVal, err := parseOptionalDate(approvedAt)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrStudentEnrollmentConditionFail)
	}

	item, err := s.db.CreateStudentEnrollment(ctx, &ent.StudentEnrollment{
		StudentID:            studentID,
		SchoolID:             schoolID,
		AcademicYearID:       academicYearID,
		ClassroomID:          classroomID,
		EnrolledAt:           enrolledAtVal,
		ExitedAt:             exitedAtVal,
		Status:               statusVal,
		EnrollmentType:       enrollmentTypeVal,
		ExitReason:           exitReasonVal,
		ExitNote:             exitNote,
		PreviousEnrollmentID: previousEnrollmentIDVal,
		RollNo:               rollNo,
		ApprovedBy:           approvedByVal,
		ApprovedAt:           approvedAtVal,
		ApprovalNote:         approvalNote,
		CreatedBy:            createdByVal,
		UpdatedBy:            updatedByVal,
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

func parseOptionalDate(v *string) (*time.Time, error) {
	if v == nil || *v == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", *v)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

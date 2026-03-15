package attendancesessions

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, schoolID uuid.UUID, academicYearID uuid.UUID, classroomID uuid.UUID, subjectID *string, teacherID *string, sessionDate string, periodNo int, mode string, startedAt *string, closedAt *string, note *string) (*ent.AttendanceSession, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "attendancesessions.service.create")
	defer span.End()

	if periodNo <= 0 {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}

	modeVal, ok := parseAttendanceMode(mode)
	if !ok {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}

	sessionDateVal, err := time.Parse("2006-01-02", sessionDate)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}

	subjectIDVal, err := parseOptionalUUID(subjectID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}
	teacherIDVal, err := parseOptionalUUID(teacherID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}
	startedAtVal, err := parseOptionalDateTime(startedAt)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}
	closedAtVal, err := parseOptionalDateTime(closedAt)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}

	startedAtPersist := time.Now()
	if startedAtVal != nil {
		startedAtPersist = *startedAtVal
	}
	if closedAtVal != nil && closedAtVal.Before(startedAtPersist) {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}

	item, err := s.db.CreateAttendanceSession(ctx, &ent.AttendanceSession{
		SchoolID:       schoolID,
		AcademicYearID: academicYearID,
		ClassroomID:    classroomID,
		SubjectID:      subjectIDVal,
		TeacherID:      teacherIDVal,
		SessionDate:    sessionDateVal,
		PeriodNo:       periodNo,
		Mode:           modeVal,
		StartedAt:      startedAtPersist,
		ClosedAt:       closedAtVal,
		Note:           note,
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

func parseOptionalDateTime(v *string) (*time.Time, error) {
	if v == nil || *v == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, *v)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseDate(v string) (time.Time, error) {
	return time.Parse("2006-01-02", v)
}

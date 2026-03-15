package attendancesessions

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.AttendanceSession, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "attendancesessions.service.update")
	defer span.End()

	if req.SchoolID == nil && req.AcademicYearID == nil && req.ClassroomID == nil && req.SubjectID == nil && req.TeacherID == nil && req.SessionDate == nil && req.PeriodNo == nil && req.Mode == nil && req.StartedAt == nil && req.ClosedAt == nil && req.Note == nil {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}

	payload := &ent.AttendanceSessionUpdate{Note: req.Note}

	var err error
	payload.SchoolID, err = parseOptionalUUID(req.SchoolID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}
	payload.AcademicYearID, err = parseOptionalUUID(req.AcademicYearID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}
	payload.ClassroomID, err = parseOptionalUUID(req.ClassroomID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}
	payload.SubjectID, err = parseOptionalUUID(req.SubjectID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}
	payload.TeacherID, err = parseOptionalUUID(req.TeacherID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}

	payload.StartedAt, err = parseOptionalDateTime(req.StartedAt)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}
	payload.ClosedAt, err = parseOptionalDateTime(req.ClosedAt)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}
	if req.SessionDate != nil && *req.SessionDate != "" {
		parsed, err := parseDate(*req.SessionDate)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
		}
		payload.SessionDate = &parsed
	}

	if req.PeriodNo != nil {
		if *req.PeriodNo <= 0 {
			return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
		}
		payload.PeriodNo = req.PeriodNo
	}

	if req.Mode != nil {
		parsed, ok := parseAttendanceMode(*req.Mode)
		if !ok {
			return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
		}
		payload.Mode = &parsed
	}

	if payload.StartedAt != nil && payload.ClosedAt != nil && payload.ClosedAt.Before(*payload.StartedAt) {
		return nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
	}

	item, err := s.db.UpdateAttendanceSessionByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

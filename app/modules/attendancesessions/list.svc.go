package attendancesessions

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, schoolID *uuid.UUID, academicYearID *uuid.UUID, classroomID *uuid.UUID, subjectID *uuid.UUID, teacherID *uuid.UUID, mode *string, sessionDateFrom *string, sessionDateTo *string) ([]*ent.AttendanceSession, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "attendancesessions.service.list")
	defer span.End()

	var modeVal *ent.AttendanceMode
	if mode != nil {
		parsed, ok := parseAttendanceMode(*mode)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrAttendanceSessionConditionFail)
		}
		modeVal = &parsed
	}

	items, page, err := s.db.ListAttendanceSessions(ctx, req, schoolID, academicYearID, classroomID, subjectID, teacherID, modeVal, sessionDateFrom, sessionDateTo)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrAttendanceSessionConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

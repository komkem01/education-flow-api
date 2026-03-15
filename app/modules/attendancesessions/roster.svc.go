package attendancesessions

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

type RosterItem struct {
	EnrollmentID uuid.UUID             `json:"enrollment_id"`
	StudentID    uuid.UUID             `json:"student_id"`
	Status       *ent.AttendanceStatus `json:"status,omitempty"`
	Source       *ent.AttendanceSource `json:"source,omitempty"`
	Remark       *string               `json:"remark,omitempty"`
}

func (s *Service) Roster(ctx context.Context, sessionID uuid.UUID) ([]*RosterItem, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "attendancesessions.service.roster")
	defer span.End()

	session, err := s.db.GetAttendanceSessionByID(ctx, sessionID)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	active := ent.StudentEnrollmentStatusActive
	enrollments, _, err := s.enrollment.ListStudentEnrollments(ctx, &base.RequestPaginate{Page: 1, Size: 10000}, nil, &session.SchoolID, &session.AcademicYearID, &session.ClassroomID, &active, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAttendanceSessionConditionFail, err)
	}

	records, _, err := s.record.ListAttendanceRecords(ctx, &base.RequestPaginate{Page: 1, Size: 10000}, &sessionID, nil, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAttendanceSessionConditionFail, err)
	}

	recMap := make(map[uuid.UUID]*ent.AttendanceRecord)
	for _, rec := range records {
		recMap[rec.EnrollmentID] = rec
	}

	items := make([]*RosterItem, 0, len(enrollments))
	for _, en := range enrollments {
		item := &RosterItem{EnrollmentID: en.ID, StudentID: en.StudentID}
		if rec, ok := recMap[en.ID]; ok {
			status := rec.Status
			source := rec.Source
			item.Status = &status
			item.Source = &source
			item.Remark = rec.Remark
		}
		items = append(items, item)
	}

	return items, nil
}

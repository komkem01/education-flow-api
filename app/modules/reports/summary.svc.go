package reports

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

type SummaryRequest struct {
	AcademicYearID *uuid.UUID
	SemesterNo     *int
}

type SummaryData struct {
	TeachersTotal     int64 `json:"teachers_total"`
	TeachersActive    int64 `json:"teachers_active"`
	StudentsTotal     int64 `json:"students_total"`
	StudentsActive    int64 `json:"students_active"`
	SubjectsTotal     int64 `json:"subjects_total"`
	CoursesTotal      int64 `json:"courses_total"`
	GradeRecordsTotal int64 `json:"grade_records_total"`
	GradePassTotal    int64 `json:"grade_pass_total"`
	GradeFailTotal    int64 `json:"grade_fail_total"`
	AttendanceTotal   int64 `json:"attendance_total"`
	AttendancePresent int64 `json:"attendance_present"`
	AttendanceAbsent  int64 `json:"attendance_absent"`
	BehaviorTotal     int64 `json:"behavior_total"`
	BehaviorGood      int64 `json:"behavior_good"`
	BehaviorBad       int64 `json:"behavior_bad"`
}

type SummaryResponse struct {
	Summary *SummaryData `json:"summary"`
}

func (s *Service) Summary(ctx context.Context, schoolID uuid.UUID, req *SummaryRequest) (*SummaryResponse, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "reports.service.summary")
	defer span.End()

	if req == nil {
		req = &SummaryRequest{}
	}

	teachersTotal, err := s.db.CountMembersByRole(ctx, schoolID, ent.MemberRoleTeacher, nil)
	if err != nil {
		return nil, fmt.Errorf("count-teachers-total: %w", err)
	}

	active := true
	teachersActive, err := s.db.CountMembersByRole(ctx, schoolID, ent.MemberRoleTeacher, &active)
	if err != nil {
		return nil, fmt.Errorf("count-teachers-active: %w", err)
	}

	studentsTotal, err := s.db.CountMembersByRole(ctx, schoolID, ent.MemberRoleStudent, nil)
	if err != nil {
		return nil, fmt.Errorf("count-students-total: %w", err)
	}

	studentsActive, err := s.db.CountMembersByRole(ctx, schoolID, ent.MemberRoleStudent, &active)
	if err != nil {
		return nil, fmt.Errorf("count-students-active: %w", err)
	}

	subjectsTotal, err := s.db.CountSubjectsBySchool(ctx, schoolID)
	if err != nil {
		return nil, fmt.Errorf("count-subjects-total: %w", err)
	}

	attendanceTotal, err := s.db.CountAttendanceByStatus(ctx, schoolID, req.AcademicYearID, nil)
	if err != nil {
		return nil, fmt.Errorf("count-attendance-total: %w", err)
	}

	present := ent.AttendanceStatusPresent
	attendancePresent, err := s.db.CountAttendanceByStatus(ctx, schoolID, req.AcademicYearID, &present)
	if err != nil {
		return nil, fmt.Errorf("count-attendance-present: %w", err)
	}

	absent := ent.AttendanceStatusAbsent
	attendanceAbsent, err := s.db.CountAttendanceByStatus(ctx, schoolID, req.AcademicYearID, &absent)
	if err != nil {
		return nil, fmt.Errorf("count-attendance-absent: %w", err)
	}

	return &SummaryResponse{
		Summary: &SummaryData{
			TeachersTotal:     teachersTotal,
			TeachersActive:    teachersActive,
			StudentsTotal:     studentsTotal,
			StudentsActive:    studentsActive,
			SubjectsTotal:     subjectsTotal,
			CoursesTotal:      0,
			GradeRecordsTotal: 0,
			GradePassTotal:    0,
			GradeFailTotal:    0,
			AttendanceTotal:   attendanceTotal,
			AttendancePresent: attendancePresent,
			AttendanceAbsent:  attendanceAbsent,
			BehaviorTotal:     0,
			BehaviorGood:      0,
			BehaviorBad:       0,
		},
	}, nil
}

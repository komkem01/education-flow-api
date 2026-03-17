package reports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (s *Service) Info(ctx context.Context, schoolID uuid.UUID, reportID string, req *SummaryRequest) (map[string]any, error) {
	summary, err := s.Summary(ctx, schoolID, req)
	if err != nil {
		return nil, err
	}

	title := "รายละเอียดรายงาน"
	description := "ข้อมูลเชิงลึกจากตัวเลขสรุป"
	metrics := map[string]any{}

	summaryData := summary.Summary
	switch reportID {
	case "teacher-summary":
		title = "สรุปข้อมูลครู"
		description = "จำนวนครูทั้งหมดและครูที่ใช้งานอยู่"
		metrics["teachers_total"] = summaryData.TeachersTotal
		metrics["teachers_active"] = summaryData.TeachersActive
	case "student-summary":
		title = "สรุปข้อมูลนักเรียน"
		description = "จำนวนนักเรียนทั้งหมดและนักเรียนที่ใช้งานอยู่"
		metrics["students_total"] = summaryData.StudentsTotal
		metrics["students_active"] = summaryData.StudentsActive
	case "grade-report":
		title = "รายงานผลการเรียน"
		description = "สรุปผลการเรียนผ่านและไม่ผ่าน"
		metrics["grade_records_total"] = summaryData.GradeRecordsTotal
		metrics["grade_pass_total"] = summaryData.GradePassTotal
		metrics["grade_fail_total"] = summaryData.GradeFailTotal
	case "attendance-report":
		title = "รายงานการเข้าเรียน"
		description = "สถิติการมาเรียนและขาดเรียน"
		metrics["attendance_total"] = summaryData.AttendanceTotal
		metrics["attendance_present"] = summaryData.AttendancePresent
		metrics["attendance_absent"] = summaryData.AttendanceAbsent
	case "behavior-report":
		title = "รายงานพฤติกรรม"
		description = "สรุปพฤติกรรมเชิงบวกและเชิงลบ"
		metrics["behavior_total"] = summaryData.BehaviorTotal
		metrics["behavior_good"] = summaryData.BehaviorGood
		metrics["behavior_bad"] = summaryData.BehaviorBad
	default:
		metrics["summary"] = summaryData
	}

	response := map[string]any{
		"id":               reportID,
		"title":            title,
		"description":      description,
		"generated_at":     time.Now().UTC(),
		"academic_year_id": nil,
		"semester_no":      nil,
		"metrics":          metrics,
	}

	if req != nil {
		if req.AcademicYearID != nil {
			response["academic_year_id"] = req.AcademicYearID.String()
		}
		if req.SemesterNo != nil {
			response["semester_no"] = *req.SemesterNo
		}
	}

	return response, nil
}

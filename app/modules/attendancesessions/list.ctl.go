package attendancesessions

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	SchoolID        string `form:"school_id"`
	AcademicYearID  string `form:"academic_year_id"`
	ClassroomID     string `form:"classroom_id"`
	SubjectID       string `form:"subject_id"`
	TeacherID       string `form:"teacher_id"`
	Mode            string `form:"mode"`
	SessionDateFrom string `form:"session_date_from"`
	SessionDateTo   string `form:"session_date_to"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	schoolID, err := parseUUIDFilter(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, "invalid-school-id", nil)
		return
	}
	academicYearID, err := parseUUIDFilter(req.AcademicYearID)
	if err != nil {
		base.BadRequest(ctx, "invalid-academic-year-id", nil)
		return
	}
	classroomID, err := parseUUIDFilter(req.ClassroomID)
	if err != nil {
		base.BadRequest(ctx, "invalid-classroom-id", nil)
		return
	}
	subjectID, err := parseUUIDFilter(req.SubjectID)
	if err != nil {
		base.BadRequest(ctx, "invalid-subject-id", nil)
		return
	}
	teacherID, err := parseUUIDFilter(req.TeacherID)
	if err != nil {
		base.BadRequest(ctx, "invalid-teacher-id", nil)
		return
	}

	var mode *string
	if req.Mode != "" {
		mode = &req.Mode
	}
	var sessionDateFrom *string
	if req.SessionDateFrom != "" {
		sessionDateFrom = &req.SessionDateFrom
	}
	var sessionDateTo *string
	if req.SessionDateTo != "" {
		sessionDateTo = &req.SessionDateTo
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, schoolID, academicYearID, classroomID, subjectID, teacherID, mode, sessionDateFrom, sessionDateTo)
	if err != nil {
		c.handleServiceError(ctx, log, err, "attendance-session-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) AttendanceSessionsList(ctx *gin.Context) {
	c.List(ctx)
}

func parseUUIDFilter(v string) (*uuid.UUID, error) {
	if v == "" {
		return nil, nil
	}
	parsed, err := uuid.Parse(v)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

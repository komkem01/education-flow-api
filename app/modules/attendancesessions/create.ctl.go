package attendancesessions

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	SchoolID       string  `json:"school_id" binding:"required"`
	AcademicYearID string  `json:"academic_year_id" binding:"required"`
	ClassroomID    string  `json:"classroom_id" binding:"required"`
	SubjectID      *string `json:"subject_id"`
	TeacherID      *string `json:"teacher_id"`
	SessionDate    string  `json:"session_date" binding:"required"`
	PeriodNo       int     `json:"period_no" binding:"required"`
	Mode           string  `json:"mode" binding:"required"`
	StartedAt      *string `json:"started_at"`
	ClosedAt       *string `json:"closed_at"`
	Note           *string `json:"note"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, "invalid-school-id", nil)
		return
	}
	academicYearID, err := uuid.Parse(req.AcademicYearID)
	if err != nil {
		base.BadRequest(ctx, "invalid-academic-year-id", nil)
		return
	}
	classroomID, err := uuid.Parse(req.ClassroomID)
	if err != nil {
		base.BadRequest(ctx, "invalid-classroom-id", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), schoolID, academicYearID, classroomID, req.SubjectID, req.TeacherID, req.SessionDate, req.PeriodNo, req.Mode, req.StartedAt, req.ClosedAt, req.Note)
	if err != nil {
		c.handleServiceError(ctx, log, err, "attendance-session-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateAttendanceSessionController(ctx *gin.Context) {
	c.Create(ctx)
}

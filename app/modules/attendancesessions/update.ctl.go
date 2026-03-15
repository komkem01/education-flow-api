package attendancesessions

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateRequest struct {
	SchoolID       *string `json:"school_id"`
	AcademicYearID *string `json:"academic_year_id"`
	ClassroomID    *string `json:"classroom_id"`
	SubjectID      *string `json:"subject_id"`
	TeacherID      *string `json:"teacher_id"`
	SessionDate    *string `json:"session_date"`
	PeriodNo       *int    `json:"period_no"`
	Mode           *string `json:"mode"`
	StartedAt      *string `json:"started_at"`
	ClosedAt       *string `json:"closed_at"`
	Note           *string `json:"note"`
}

func (c *Controller) Update(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var uriReq UpdateByIDRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	id, err := uuid.Parse(uriReq.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return
	}

	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	item, err := c.svc.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		c.handleServiceError(ctx, log, err, "attendance-session-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) AttendanceSessionsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

package attendancerecords

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	SessionID    string  `json:"session_id" binding:"required"`
	EnrollmentID string  `json:"enrollment_id" binding:"required"`
	Status       string  `json:"status" binding:"required"`
	Source       string  `json:"source" binding:"required"`
	MarkedAt     *string `json:"marked_at"`
	Remark       *string `json:"remark"`
	MarkedBy     *string `json:"marked_by"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		base.BadRequest(ctx, "invalid-session-id", nil)
		return
	}
	enrollmentID, err := uuid.Parse(req.EnrollmentID)
	if err != nil {
		base.BadRequest(ctx, "invalid-enrollment-id", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), sessionID, enrollmentID, req.Status, req.Source, req.MarkedAt, req.Remark, req.MarkedBy)
	if err != nil {
		c.handleServiceError(ctx, log, err, "attendance-record-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateAttendanceRecordController(ctx *gin.Context) {
	c.Create(ctx)
}

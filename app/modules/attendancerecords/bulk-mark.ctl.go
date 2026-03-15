package attendancerecords

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

type BulkMarkItem struct {
	EnrollmentID string  `json:"enrollment_id" binding:"required"`
	Status       string  `json:"status" binding:"required"`
	Remark       *string `json:"remark"`
}

type BulkMarkRequest struct {
	SessionID string         `json:"session_id" binding:"required"`
	Source    string         `json:"source" binding:"required"`
	MarkedBy  *string        `json:"marked_by"`
	Reason    *string        `json:"reason"`
	Items     []BulkMarkItem `json:"items" binding:"required"`
}

func (c *Controller) BulkMark(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req BulkMarkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	if len(req.Items) == 0 {
		base.BadRequest(ctx, "items-required", nil)
		return
	}

	result, err := c.svc.BulkMark(ctx.Request.Context(), &req)
	if err != nil {
		c.handleServiceError(ctx, log, err, "attendance-record-bulk-mark-failed")
		return
	}

	base.Success(ctx, result, "success")
}

func (c *Controller) AttendanceRecordsBulkMark(ctx *gin.Context) {
	c.BulkMark(ctx)
}

package attendancerecordlogs

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	RecordID  string  `json:"record_id" binding:"required"`
	OldStatus *string `json:"old_status"`
	NewStatus string  `json:"new_status" binding:"required"`
	ChangedBy *string `json:"changed_by"`
	ChangedAt *string `json:"changed_at"`
	Reason    *string `json:"reason"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	recordID, err := uuid.Parse(req.RecordID)
	if err != nil {
		base.BadRequest(ctx, "invalid-record-id", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), recordID, req.OldStatus, req.NewStatus, req.ChangedBy, req.ChangedAt, req.Reason)
	if err != nil {
		c.handleServiceError(ctx, log, err, "attendance-record-log-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateAttendanceRecordLogController(ctx *gin.Context) {
	c.Create(ctx)
}

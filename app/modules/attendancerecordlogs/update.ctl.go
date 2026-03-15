package attendancerecordlogs

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
	RecordID  *string `json:"record_id"`
	OldStatus *string `json:"old_status"`
	NewStatus *string `json:"new_status"`
	ChangedBy *string `json:"changed_by"`
	ChangedAt *string `json:"changed_at"`
	Reason    *string `json:"reason"`
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
		c.handleServiceError(ctx, log, err, "attendance-record-log-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) AttendanceRecordLogsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

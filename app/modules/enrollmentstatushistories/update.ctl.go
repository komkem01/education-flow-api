package enrollmentstatushistories

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
	EnrollmentID *string `json:"enrollment_id"`
	FromStatus   *string `json:"from_status"`
	ToStatus     *string `json:"to_status"`
	ChangedAt    *string `json:"changed_at"`
	ChangedBy    *string `json:"changed_by"`
	Reason       *string `json:"reason"`
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
		c.handleServiceError(ctx, log, err, "enrollment-status-history-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) EnrollmentStatusHistoriesUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

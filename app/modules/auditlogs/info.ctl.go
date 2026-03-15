package auditlogs

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) GetByID(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req GetByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		c.handleServiceError(ctx, log, err, "audit-log-get-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) AuditLogsInfo(ctx *gin.Context) {
	c.GetByID(ctx)
}

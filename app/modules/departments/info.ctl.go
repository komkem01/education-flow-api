package departments

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
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return
	}

	department, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		c.handleServiceError(ctx, log, err, "department-get-failed")
		return
	}

	base.Success(ctx, department, "success")
}

func (c *Controller) InfoController(ctx *gin.Context) {
	c.GetByID(ctx)
}

func (c *Controller) DepartmentsInfo(ctx *gin.Context) {
	c.GetByID(ctx)
}

package departments

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateRequest struct {
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
}

type UpdateByIDRequest struct {
	ID string `uri:"id" binding:"required"`
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

	department, err := c.svc.Update(ctx.Request.Context(), id, req.Name, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "department-update-failed")
		return
	}

	base.Success(ctx, department, "success")
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	c.Update(ctx)
}

func (c *Controller) DepartmentsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

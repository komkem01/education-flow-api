package schools

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateRequest struct {
	Name        *string `json:"name"`
	LogoURL     *string `json:"logo_url"`
	ThemeColor  *string `json:"theme_color"`
	Address     *string `json:"address"`
	Description *string `json:"description"`
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

	school, err := c.svc.Update(ctx.Request.Context(), id, req.Name, req.LogoURL, req.ThemeColor, req.Address, req.Description)
	if err != nil {
		c.handleServiceError(ctx, log, err, "school-update-failed")
		return
	}

	base.Success(ctx, school, "success")
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	c.Update(ctx)
}

func (c *Controller) SchoolsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

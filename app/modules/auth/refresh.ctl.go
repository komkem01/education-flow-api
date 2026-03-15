package auth

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (c *Controller) Refresh(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.Refresh(ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		c.handleServiceError(ctx, log, err, "auth-refresh-failed")
		return
	}

	base.Success(ctx, res, "success")
}

func (c *Controller) AuthRefresh(ctx *gin.Context) {
	c.Refresh(ctx)
}

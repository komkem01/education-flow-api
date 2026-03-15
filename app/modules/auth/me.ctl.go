package auth

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

func (c *Controller) Me(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	res, err := c.svc.Me(ctx.Request.Context(), ctx.GetHeader("Authorization"))
	if err != nil {
		c.handleServiceError(ctx, log, err, "auth-me-failed")
		return
	}

	base.Success(ctx, res, "success")
}

func (c *Controller) AuthMe(ctx *gin.Context) {
	c.Me(ctx)
}

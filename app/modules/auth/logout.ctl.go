package auth

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

func (c *Controller) Logout(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	if err := c.svc.Logout(ctx.Request.Context(), ctx.GetHeader("Authorization")); err != nil {
		c.handleServiceError(ctx, log, err, "auth-logout-failed")
		return
	}

	base.Success(ctx, nil, "success")
}

func (c *Controller) LogoutAll(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	if err := c.svc.LogoutAll(ctx.Request.Context(), ctx.GetHeader("Authorization")); err != nil {
		c.handleServiceError(ctx, log, err, "auth-logout-all-failed")
		return
	}

	base.Success(ctx, nil, "success")
}

func (c *Controller) AuthLogout(ctx *gin.Context) {
	c.Logout(ctx)
}

func (c *Controller) AuthLogoutAll(ctx *gin.Context) {
	c.LogoutAll(ctx)
}

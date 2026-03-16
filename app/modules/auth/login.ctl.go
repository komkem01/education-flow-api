package auth

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (c *Controller) Login(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	res, err := c.svc.Login(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.handleServiceError(ctx, log, err, "auth-login-failed")
		return
	}

	c.writeAuthCookies(ctx, res.Token, res.RefreshToken)

	base.Success(ctx, res, "success")
}

func (c *Controller) AuthLogin(ctx *gin.Context) {
	c.Login(ctx)
}

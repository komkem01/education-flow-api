package auth

import (
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

const CurrentUserContextKey = "current_user"

func (c *Controller) RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := c.svc.resolveCurrentUser(ctx.Request.Context(), ctx.GetHeader("Authorization"))
		if err != nil {
			base.Unauthorized(ctx, "unauthorized", nil)
			ctx.Abort()
			return
		}

		ctx.Set(CurrentUserContextKey, user)
		ctx.Next()
	}
}

func CurrentUserFromGin(ctx *gin.Context) (*CurrentUser, bool) {
	v, ok := ctx.Get(CurrentUserContextKey)
	if !ok {
		return nil, false
	}
	user, ok := v.(*CurrentUser)
	if !ok || user == nil {
		return nil, false
	}
	return user, true
}

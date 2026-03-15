package storages

import (
	"eduflow/app/modules/auth"
	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Provider   string         `json:"provider" binding:"required"`
	Name       string         `json:"name" binding:"required"`
	Endpoint   *string        `json:"endpoint"`
	BucketName string         `json:"bucket_name" binding:"required"`
	IsDefault  *bool          `json:"is_default"`
	Config     map[string]any `json:"config"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	currentUser, ok := auth.CurrentUserFromGin(ctx)
	if !ok || currentUser == nil || currentUser.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}
	if currentUser.Member.Role != ent.MemberRoleAdmin && currentUser.Member.Role != ent.MemberRoleSuperadmin {
		base.Forbidden(ctx, "unauthorized", nil)
		return
	}

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	item, err := c.svc.Create(
		ctx.Request.Context(),
		currentUser.Member.SchoolID,
		req.Provider,
		req.Name,
		req.Endpoint,
		req.BucketName,
		req.IsDefault,
		req.Config,
	)
	if err != nil {
		c.handleServiceError(ctx, log, err, "storage-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) StoragesCreate(ctx *gin.Context) {
	c.Create(ctx)
}

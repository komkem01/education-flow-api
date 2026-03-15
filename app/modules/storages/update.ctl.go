package storages

import (
	"eduflow/app/modules/auth"
	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateRequest struct {
	Provider   *string         `json:"provider"`
	Name       *string         `json:"name"`
	Endpoint   *string         `json:"endpoint"`
	BucketName *string         `json:"bucket_name"`
	IsDefault  *bool           `json:"is_default"`
	Config     *map[string]any `json:"config"`
}

type UpdateByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) Update(ctx *gin.Context) {
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

	var uriReq UpdateByIDRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	id, err := uuid.Parse(uriReq.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	item, err := c.svc.Update(
		ctx.Request.Context(),
		id,
		currentUser.Member.SchoolID,
		req.Provider,
		req.Name,
		req.Endpoint,
		req.BucketName,
		req.IsDefault,
		req.Config,
	)
	if err != nil {
		c.handleServiceError(ctx, log, err, "storage-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) StoragesUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

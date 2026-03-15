package storages

import (
	"eduflow/app/modules/auth"
	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

type ListRequest struct {
	base.RequestPaginate
	Provider  string `form:"provider"`
	IsDefault *bool  `form:"is_default"`
}

func (c *Controller) List(ctx *gin.Context) {
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

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	var provider *string
	if req.Provider != "" {
		provider = &req.Provider
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, currentUser.Member.SchoolID, provider, req.IsDefault)
	if err != nil {
		c.handleServiceError(ctx, log, err, "storage-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) StoragesList(ctx *gin.Context) {
	c.List(ctx)
}

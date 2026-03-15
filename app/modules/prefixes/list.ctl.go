package prefixes

import (
	"strconv"

	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	IsActive *bool  `form:"is_active"`
	GenderID string `form:"gender_id"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	if raw := ctx.Query("is_active"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			base.BadRequest(ctx, "invalid-is-active", nil)
			return
		}
		req.IsActive = &parsed
	}

	var genderID *uuid.UUID
	if req.GenderID != "" {
		parsed, err := uuid.Parse(req.GenderID)
		if err != nil {
			base.BadRequest(ctx, "invalid-id", nil)
			return
		}
		genderID = &parsed
	}

	prefixes, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, req.IsActive, genderID)
	if err != nil {
		c.handleServiceError(ctx, log, err, "prefix-list-failed")
		return
	}

	base.Paginate(ctx, prefixes, page)
}

func (c *Controller) PrefixesList(ctx *gin.Context) {
	c.List(ctx)
}

package members

import (
	"strconv"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	IsActive *bool  `form:"is_active"`
	SchoolID string `form:"school_id"`
	Role     string `form:"role"`
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

	var schoolID *uuid.UUID
	if req.SchoolID != "" {
		parsed, err := uuid.Parse(req.SchoolID)
		if err != nil {
			base.BadRequest(ctx, "invalid-id", nil)
			return
		}
		schoolID = &parsed
	}

	var role *ent.MemberRole
	if req.Role != "" {
		parsed, ok := parseMemberRole(req.Role)
		if !ok {
			base.BadRequest(ctx, "invalid-role", nil)
			return
		}
		role = &parsed
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, req.IsActive, schoolID, role)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-list-failed")
		return
	}

	resp := make([]*memberResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, toMemberResponse(item))
	}

	base.Paginate(ctx, resp, page)
}

func (c *Controller) MembersList(ctx *gin.Context) {
	c.List(ctx)
}

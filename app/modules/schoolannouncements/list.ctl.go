package schoolannouncements

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
	SchoolID   string `form:"school_id"`
	Status     string `form:"status"`
	TargetRole string `form:"target_role"`
	IsPinned   *bool  `form:"is_pinned"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	user, ok := c.currentUser(ctx)
	if !ok {
		return
	}

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	if raw := ctx.Query("is_pinned"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			base.BadRequest(ctx, "invalid-is-pinned", nil)
			return
		}
		req.IsPinned = &parsed
	}

	var schoolID *uuid.UUID
	if req.SchoolID != "" {
		parsed, err := uuid.Parse(req.SchoolID)
		if err != nil {
			base.BadRequest(ctx, "invalid-school-id", nil)
			return
		}
		schoolID = &parsed
	}

	if user.Member.Role != ent.MemberRoleSuperadmin {
		schoolID = &user.Member.SchoolID
	}

	var status *ent.SchoolAnnouncementStatus
	if req.Status != "" {
		parsed := ent.SchoolAnnouncementStatus(req.Status)
		if !isValidStatus(parsed) {
			base.BadRequest(ctx, "invalid-status", nil)
			return
		}
		status = &parsed
	}

	var targetRole *string
	if req.TargetRole != "" {
		if !isValidTargetRole(req.TargetRole) {
			base.BadRequest(ctx, "invalid-target-role", nil)
			return
		}
		targetRole = &req.TargetRole
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, schoolID, status, targetRole, req.IsPinned)
	if err != nil {
		c.handleServiceError(ctx, log, err, "school-announcement-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) SchoolAnnouncementsList(ctx *gin.Context) {
	c.List(ctx)
}

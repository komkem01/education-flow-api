package teacherexperiences

import (
	"strconv"

	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	TeacherID string `form:"teacher_id"`
	IsCurrent *bool  `form:"is_current"`
	IsActive  *bool  `form:"is_active"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	if raw := ctx.Query("is_current"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			base.BadRequest(ctx, "invalid-is-current", nil)
			return
		}
		req.IsCurrent = &parsed
	}

	if raw := ctx.Query("is_active"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			base.BadRequest(ctx, "invalid-is-active", nil)
			return
		}
		req.IsActive = &parsed
	}

	var teacherID *uuid.UUID
	if req.TeacherID != "" {
		parsed, err := uuid.Parse(req.TeacherID)
		if err != nil {
			base.BadRequest(ctx, "invalid-teacher-id", nil)
			return
		}
		teacherID = &parsed
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, teacherID, req.IsCurrent, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "teacher-experience-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) TeacherExperiencesList(ctx *gin.Context) {
	c.List(ctx)
}

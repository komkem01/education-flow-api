package memberstudents

import (
	"strconv"

	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	IsActive         *bool  `form:"is_active"`
	SchoolID         string `form:"school_id"`
	AdvisorTeacherID string `form:"advisor_teacher_id"`
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
			base.BadRequest(ctx, "invalid-school-id", nil)
			return
		}
		schoolID = &parsed
	}

	var advisorTeacherID *uuid.UUID
	if req.AdvisorTeacherID != "" {
		parsed, err := uuid.Parse(req.AdvisorTeacherID)
		if err != nil {
			base.BadRequest(ctx, "invalid-advisor-teacher-id", nil)
			return
		}
		advisorTeacherID = &parsed
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, req.IsActive, schoolID, advisorTeacherID)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-student-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) MemberStudentsList(ctx *gin.Context) {
	c.List(ctx)
}

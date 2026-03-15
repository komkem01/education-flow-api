package studentguardians

import (
	"strconv"

	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	StudentID      string `form:"student_id"`
	GuardianID     string `form:"guardian_id"`
	IsMainGuardian *bool  `form:"is_main_guardian"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	if raw := ctx.Query("is_main_guardian"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			base.BadRequest(ctx, "invalid-is-main-guardian", nil)
			return
		}
		req.IsMainGuardian = &parsed
	}

	var studentID *uuid.UUID
	if req.StudentID != "" {
		parsed, err := uuid.Parse(req.StudentID)
		if err != nil {
			base.BadRequest(ctx, "invalid-student-id", nil)
			return
		}
		studentID = &parsed
	}

	var guardianID *uuid.UUID
	if req.GuardianID != "" {
		parsed, err := uuid.Parse(req.GuardianID)
		if err != nil {
			base.BadRequest(ctx, "invalid-guardian-id", nil)
			return
		}
		guardianID = &parsed
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, studentID, guardianID, req.IsMainGuardian)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-guardian-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) StudentGuardiansList(ctx *gin.Context) {
	c.List(ctx)
}

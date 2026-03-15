package teachersubjects

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
	SubjectID string `form:"subject_id"`
	Role      string `form:"role"`
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

	var subjectID *uuid.UUID
	if req.SubjectID != "" {
		parsed, err := uuid.Parse(req.SubjectID)
		if err != nil {
			base.BadRequest(ctx, "invalid-subject-id", nil)
			return
		}
		subjectID = &parsed
	}

	var role *string
	if req.Role != "" {
		role = &req.Role
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, teacherID, subjectID, role, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "teacher-subject-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) TeacherSubjectsList(ctx *gin.Context) {
	c.List(ctx)
}

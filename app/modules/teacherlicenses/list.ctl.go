package teacherlicenses

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	TeacherID string `form:"teacher_id"`
	Status    string `form:"status"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
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

	var status *string
	if req.Status != "" {
		status = &req.Status
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, teacherID, status)
	if err != nil {
		c.handleServiceError(ctx, log, err, "teacher-license-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) TeacherLicensesList(ctx *gin.Context) {
	c.List(ctx)
}

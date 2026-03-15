package teacherrequests

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	TeacherID     string `form:"teacher_id"`
	RequestType   string `form:"request_type"`
	RequestStatus string `form:"status"`
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

	var requestType *string
	if req.RequestType != "" {
		requestType = &req.RequestType
	}

	var status *string
	if req.RequestStatus != "" {
		status = &req.RequestStatus
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, teacherID, requestType, status)
	if err != nil {
		c.handleServiceError(ctx, log, err, "teacher-request-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) TeacherRequestsList(ctx *gin.Context) {
	c.List(ctx)
}

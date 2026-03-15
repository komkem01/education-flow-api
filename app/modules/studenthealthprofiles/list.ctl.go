package studenthealthprofiles

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	StudentID string `form:"student_id"`
	BloodType string `form:"blood_type"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
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

	var bloodType *string
	if req.BloodType != "" {
		bloodType = &req.BloodType
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, studentID, bloodType)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-health-profile-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) StudentHealthProfilesList(ctx *gin.Context) {
	c.List(ctx)
}

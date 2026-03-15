package teacherhealthprofiles

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	MemberTeacherID string `form:"member_teacher_id"`
	BloodType       string `form:"blood_type"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var memberTeacherID *uuid.UUID
	if req.MemberTeacherID != "" {
		parsed, err := uuid.Parse(req.MemberTeacherID)
		if err != nil {
			base.BadRequest(ctx, "invalid-member-teacher-id", nil)
			return
		}
		memberTeacherID = &parsed
	}

	var bloodType *string
	if req.BloodType != "" {
		bloodType = &req.BloodType
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, memberTeacherID, bloodType)
	if err != nil {
		c.handleServiceError(ctx, log, err, "teacher-health-profile-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) TeacherHealthProfilesList(ctx *gin.Context) {
	c.List(ctx)
}

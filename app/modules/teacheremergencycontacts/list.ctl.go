package teacheremergencycontacts

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	MemberTeacherID string `form:"member_teacher_id"`
	IsPrimary       string `form:"is_primary"`
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

	var isPrimary *bool
	if req.IsPrimary != "" {
		parsed, ok := parseBoolFilter(req.IsPrimary)
		if !ok {
			base.BadRequest(ctx, "invalid-is-primary", nil)
			return
		}
		isPrimary = &parsed
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, memberTeacherID, isPrimary)
	if err != nil {
		c.handleServiceError(ctx, log, err, "teacher-emergency-contact-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) TeacherEmergencyContactsList(ctx *gin.Context) {
	c.List(ctx)
}

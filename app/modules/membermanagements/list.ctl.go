package membermanagements

import (
	"strconv"

	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	IsActive           *bool  `form:"is_active"`
	MemberID           string `form:"member_id"`
	DepartmentID       string `form:"department_id"`
	SchoolDepartmentID string `form:"school_department_id"`
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

	var memberID *uuid.UUID
	if req.MemberID != "" {
		parsed, err := uuid.Parse(req.MemberID)
		if err != nil {
			base.BadRequest(ctx, "invalid-member-id", nil)
			return
		}
		memberID = &parsed
	}

	var departmentID *uuid.UUID
	if req.DepartmentID != "" {
		parsed, err := uuid.Parse(req.DepartmentID)
		if err != nil {
			base.BadRequest(ctx, "invalid-department-id", nil)
			return
		}
		departmentID = &parsed
	}

	var schoolDepartmentID *uuid.UUID
	if req.SchoolDepartmentID != "" {
		parsed, err := uuid.Parse(req.SchoolDepartmentID)
		if err != nil {
			base.BadRequest(ctx, "invalid-school-department-id", nil)
			return
		}
		schoolDepartmentID = &parsed
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, req.IsActive, memberID, departmentID, schoolDepartmentID)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-management-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) MemberManagementsList(ctx *gin.Context) {
	c.List(ctx)
}

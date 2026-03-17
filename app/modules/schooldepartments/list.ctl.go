package schooldepartments

import (
	"strconv"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	SchoolID     string `form:"school_id"`
	DepartmentID string `form:"department_id"`
	IsActive     *bool  `form:"is_active"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	user, ok := c.currentUser(ctx)
	if !ok {
		return
	}

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
	if user.Member.Role != ent.MemberRoleSuperadmin {
		schoolID = &user.Member.SchoolID
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

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, schoolID, departmentID, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "school-department-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) SchoolDepartmentsList(ctx *gin.Context) {
	c.List(ctx)
}

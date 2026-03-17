package schooldepartments

import (
	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateRequest struct {
	SchoolID     *string `json:"school_id"`
	DepartmentID *string `json:"department_id"`
	Code         *string `json:"code"`
	CustomName   *string `json:"custom_name"`
	IsActive     *bool   `json:"is_active"`
}

func (c *Controller) Update(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	user, ok := c.currentUser(ctx)
	if !ok {
		return
	}
	if !canWriteSchoolDepartments(user.Member.Role) {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var uriReq UpdateByIDRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	id, err := uuid.Parse(uriReq.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return
	}

	current, err := c.svc.Info(ctx.Request.Context(), id)
	if err != nil {
		c.handleServiceError(ctx, log, err, "school-department-info-failed")
		return
	}
	if user.Member.Role != ent.MemberRoleSuperadmin && current.SchoolID != user.Member.SchoolID {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	if req.Code != nil {
		base.ValidateFailed(ctx, "school-department-code-readonly", nil)
		return
	}

	payload := &ent.SchoolDepartmentUpdate{CustomName: req.CustomName, IsActive: req.IsActive}
	if req.SchoolID != nil {
		parsed, err := uuid.Parse(*req.SchoolID)
		if err != nil {
			base.BadRequest(ctx, "invalid-school-id", nil)
			return
		}
		payload.SchoolID = &parsed
	}
	if req.DepartmentID != nil {
		parsed, err := uuid.Parse(*req.DepartmentID)
		if err != nil {
			base.BadRequest(ctx, "invalid-department-id", nil)
			return
		}
		payload.DepartmentID = &parsed
	}

	item, err := c.svc.Update(ctx.Request.Context(), id, payload)
	if err != nil {
		c.handleServiceError(ctx, log, err, "school-department-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) SchoolDepartmentsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

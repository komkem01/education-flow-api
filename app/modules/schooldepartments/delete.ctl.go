package schooldepartments

import (
	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeleteByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) Delete(ctx *gin.Context) {
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

	var req DeleteByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	id, err := uuid.Parse(req.ID)
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

	if err := c.svc.Delete(ctx.Request.Context(), id); err != nil {
		c.handleServiceError(ctx, log, err, "school-department-delete-failed")
		return
	}

	base.Success(ctx, gin.H{"deleted": true}, "success")
}

func (c *Controller) SchoolDepartmentsDelete(ctx *gin.Context) {
	c.Delete(ctx)
}

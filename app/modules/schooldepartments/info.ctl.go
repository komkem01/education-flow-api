package schooldepartments

import (
	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InfoByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) Info(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	user, ok := c.currentUser(ctx)
	if !ok {
		return
	}

	var req InfoByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return
	}

	item, err := c.svc.Info(ctx.Request.Context(), id)
	if err != nil {
		c.handleServiceError(ctx, log, err, "school-department-info-failed")
		return
	}

	if user.Member.Role != ent.MemberRoleSuperadmin && item.SchoolID != user.Member.SchoolID {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) SchoolDepartmentsInfo(ctx *gin.Context) {
	c.Info(ctx)
}

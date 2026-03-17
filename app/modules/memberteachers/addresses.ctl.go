package memberteachers

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListAddressesRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) ListAddresses(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListAddressesRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return
	}

	items, err := c.svc.ListAddresses(ctx.Request.Context(), id)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-teacher-address-list-failed")
		return
	}

	base.Success(ctx, items, "success")
}

func (c *Controller) MemberTeacherAddressesList(ctx *gin.Context) {
	c.ListAddresses(ctx)
}

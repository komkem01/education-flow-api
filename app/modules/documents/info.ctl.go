package documents

import (
	"eduflow/app/modules/auth"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) GetByID(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	currentUser, ok := auth.CurrentUserFromGin(ctx)
	if !ok || currentUser == nil || currentUser.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var req GetByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	item, err := c.svc.GetByID(ctx.Request.Context(), id, currentUser.Member.SchoolID)
	if err != nil {
		c.handleServiceError(ctx, log, err, "document-get-failed")
		return
	}
	if !c.isElevatedRole(currentUser.Member.Role) && (item.OwnerMemberID == nil || *item.OwnerMemberID != currentUser.Member.ID) {
		base.Forbidden(ctx, "unauthorized", nil)
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) DocumentsInfo(ctx *gin.Context) {
	c.GetByID(ctx)
}

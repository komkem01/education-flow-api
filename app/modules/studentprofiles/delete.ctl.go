package studentprofiles

import (
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

	if err := c.svc.Delete(ctx.Request.Context(), id); err != nil {
		c.handleServiceError(ctx, log, err, "student-profile-delete-failed")
		return
	}

	base.Success(ctx, gin.H{"deleted": true}, "success")
}

func (c *Controller) StudentProfilesDelete(ctx *gin.Context) {
	c.Delete(ctx)
}

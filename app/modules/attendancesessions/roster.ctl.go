package attendancesessions

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (c *Controller) Roster(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	idRaw := ctx.Param("id")
	sessionID, err := uuid.Parse(idRaw)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return
	}

	items, err := c.svc.Roster(ctx.Request.Context(), sessionID)
	if err != nil {
		c.handleServiceError(ctx, log, err, "attendance-session-roster-failed")
		return
	}

	base.Success(ctx, items, "success")
}

func (c *Controller) AttendanceSessionsRoster(ctx *gin.Context) {
	c.Roster(ctx)
}

package reports

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

func (c *Controller) Filters(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	user, ok := c.currentUser(ctx)
	if !ok {
		return
	}

	if !c.canRead(user) {
		base.Forbidden(ctx, "forbidden", nil)
		return
	}

	res, err := c.svc.Filters(ctx.Request.Context(), fallbackSchoolID(user.Member))
	if err != nil {
		log.Errf("report-filters-failed: %v", err)
		base.InternalServerError(ctx, "report-filters-failed", nil)
		return
	}

	base.Success(ctx, res, "success")
}

func (c *Controller) ReportsFilters(ctx *gin.Context) {
	c.Filters(ctx)
}

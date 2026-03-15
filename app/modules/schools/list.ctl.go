package schools

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

type ListRequest struct {
	base.RequestPaginate
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	schools, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate)
	if err != nil {
		c.handleServiceError(ctx, log, err, "school-list-failed")
		return
	}

	base.Paginate(ctx, schools, page)
}

func (c *Controller) SchoolsList(ctx *gin.Context) {
	c.List(ctx)
}

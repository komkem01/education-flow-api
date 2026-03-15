package departments

import (
	"strconv"

	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

type ListRequest struct {
	base.RequestPaginate
	IsActive *bool `form:"is_active"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if raw := ctx.Query("is_active"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			base.BadRequest(ctx, "invalid-is-active", nil)
			return
		}
		req.IsActive = &parsed
	}

	departments, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "department-list-failed")
		return
	}

	base.Paginate(ctx, departments, page)
}

func (c *Controller) DepartmentsList(ctx *gin.Context) {
	c.List(ctx)
}

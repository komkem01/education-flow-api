package departments

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	IsActive bool   `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	department, err := c.svc.Create(ctx.Request.Context(), req.Code, req.Name, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "department-create-failed")
		return
	}

	base.Success(ctx, department, "success")
}

func (c *Controller) CreateDepartmentController(ctx *gin.Context) {
	c.Create(ctx)
}

func (c *Controller) DepartmentsCreate(ctx *gin.Context) {
	c.Create(ctx)
}

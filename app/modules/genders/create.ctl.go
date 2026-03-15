package genders

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
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

	gender, err := c.svc.Create(ctx.Request.Context(), req.Name, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "gender-create-failed")
		return
	}

	base.Success(ctx, gender, "success")
}

func (c *Controller) CreateGenderController(ctx *gin.Context) {
	c.Create(ctx)
}

func (c *Controller) GendersCreate(ctx *gin.Context) {
	c.Create(ctx)
}

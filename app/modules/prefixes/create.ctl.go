package prefixes

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	GenderID string `json:"gender_id" binding:"required"`
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

	genderID, err := uuid.Parse(req.GenderID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return
	}

	prefix, err := c.svc.Create(ctx.Request.Context(), genderID, req.Name, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "prefix-create-failed")
		return
	}

	base.Success(ctx, prefix, "success")
}

func (c *Controller) CreatePrefixController(ctx *gin.Context) {
	c.Create(ctx)
}

func (c *Controller) PrefixesCreate(ctx *gin.Context) {
	c.Create(ctx)
}

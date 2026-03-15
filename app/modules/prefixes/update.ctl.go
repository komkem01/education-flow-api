package prefixes

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateRequest struct {
	GenderID *string `json:"gender_id"`
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
}

type UpdateByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) Update(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var uriReq UpdateByIDRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	id, err := uuid.Parse(uriReq.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return
	}

	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var genderID *uuid.UUID
	if req.GenderID != nil {
		parsed, err := uuid.Parse(*req.GenderID)
		if err != nil {
			base.BadRequest(ctx, "invalid-id", nil)
			return
		}
		genderID = &parsed
	}

	prefix, err := c.svc.Update(ctx.Request.Context(), id, genderID, req.Name, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "prefix-update-failed")
		return
	}

	base.Success(ctx, prefix, "success")
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	c.Update(ctx)
}

func (c *Controller) PrefixesUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

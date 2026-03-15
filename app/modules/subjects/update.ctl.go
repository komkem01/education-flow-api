package subjects

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateRequest struct {
	SchoolID       *string  `json:"school_id"`
	SubjectGroupID *string  `json:"subject_group_id"`
	Code           *string  `json:"code"`
	NameTH         *string  `json:"name_th"`
	NameEN         *string  `json:"name_en"`
	Credit         *float64 `json:"credit"`
	HoursPerWeek   *int     `json:"hours_per_week"`
	IsElective     *bool    `json:"is_elective"`
	IsActive       *bool    `json:"is_active"`
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

	item, err := c.svc.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		c.handleServiceError(ctx, log, err, "subject-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) SubjectsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

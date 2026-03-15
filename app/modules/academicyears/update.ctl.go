package academicyears

import (
	"time"

	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateRequest struct {
	SchoolID  *string `json:"school_id"`
	Year      *string `json:"year"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	IsActive  *bool   `json:"is_active"`
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

	var schoolID *uuid.UUID
	if req.SchoolID != nil {
		parsed, err := uuid.Parse(*req.SchoolID)
		if err != nil {
			base.BadRequest(ctx, "invalid-id", nil)
			return
		}
		schoolID = &parsed
	}

	var startDate *time.Time
	if req.StartDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			base.BadRequest(ctx, "invalid-start-date", nil)
			return
		}
		startDate = &parsed
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			base.BadRequest(ctx, "invalid-end-date", nil)
			return
		}
		endDate = &parsed
	}

	item, err := c.svc.Update(ctx.Request.Context(), id, schoolID, req.Year, startDate, endDate, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "academic-year-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	c.Update(ctx)
}

func (c *Controller) AcademicYearsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

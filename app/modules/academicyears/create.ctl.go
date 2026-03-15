package academicyears

import (
	"time"

	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	SchoolID  string `json:"school_id" binding:"required"`
	Year      string `json:"year" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	IsActive  bool   `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		base.BadRequest(ctx, "invalid-start-date", nil)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		base.BadRequest(ctx, "invalid-end-date", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), schoolID, req.Year, startDate, endDate, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "academic-year-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateAcademicYearController(ctx *gin.Context) {
	c.Create(ctx)
}

func (c *Controller) AcademicYearsCreate(ctx *gin.Context) {
	c.Create(ctx)
}

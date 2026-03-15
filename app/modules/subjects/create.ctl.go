package subjects

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	SchoolID       string  `json:"school_id" binding:"required"`
	SubjectGroupID string  `json:"subject_group_id" binding:"required"`
	Code           string  `json:"code" binding:"required"`
	NameTH         string  `json:"name_th" binding:"required"`
	NameEN         *string `json:"name_en"`
	Credit         float64 `json:"credit"`
	HoursPerWeek   *int    `json:"hours_per_week"`
	IsElective     bool    `json:"is_elective"`
	IsActive       bool    `json:"is_active"`
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
		base.BadRequest(ctx, "invalid-school-id", nil)
		return
	}
	subjectGroupID, err := uuid.Parse(req.SubjectGroupID)
	if err != nil {
		base.BadRequest(ctx, "invalid-subject-group-id", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), schoolID, subjectGroupID, req.Code, req.NameTH, req.NameEN, req.Credit, req.HoursPerWeek, req.IsElective, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "subject-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateSubjectController(ctx *gin.Context) {
	c.Create(ctx)
}

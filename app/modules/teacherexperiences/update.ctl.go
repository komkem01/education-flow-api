package teacherexperiences

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
	TeacherID        *string `json:"teacher_id"`
	SchoolName       *string `json:"school_name"`
	Position         *string `json:"position"`
	DepartmentName   *string `json:"department_name"`
	StartDate        *string `json:"start_date"`
	EndDate          *string `json:"end_date"`
	IsCurrent        *bool   `json:"is_current"`
	Responsibilities *string `json:"responsibilities"`
	Achievements     *string `json:"achievements"`
	SortOrder        *int    `json:"sort_order"`
	IsActive         *bool   `json:"is_active"`
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
		c.handleServiceError(ctx, log, err, "teacher-experience-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) TeacherExperiencesUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

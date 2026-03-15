package teacherexperiences

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	TeacherID        string  `json:"teacher_id" binding:"required"`
	SchoolName       string  `json:"school_name" binding:"required"`
	Position         string  `json:"position" binding:"required"`
	DepartmentName   *string `json:"department_name"`
	StartDate        string  `json:"start_date" binding:"required"`
	EndDate          *string `json:"end_date"`
	IsCurrent        bool    `json:"is_current"`
	Responsibilities *string `json:"responsibilities"`
	Achievements     *string `json:"achievements"`
	SortOrder        int     `json:"sort_order"`
	IsActive         bool    `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	teacherID, err := uuid.Parse(req.TeacherID)
	if err != nil {
		base.BadRequest(ctx, "invalid-teacher-id", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), teacherID, req.SchoolName, req.Position, req.DepartmentName, req.StartDate, req.EndDate, req.IsCurrent, req.Responsibilities, req.Achievements, req.SortOrder, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "teacher-experience-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateTeacherExperienceController(ctx *gin.Context) {
	c.Create(ctx)
}

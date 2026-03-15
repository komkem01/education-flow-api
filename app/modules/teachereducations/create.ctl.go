package teachereducations

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	TeacherID      string `json:"teacher_id" binding:"required"`
	Degree         string `json:"degree" binding:"required"`
	Major          string `json:"major" binding:"required"`
	University     string `json:"university" binding:"required"`
	GraduationYear string `json:"graduation_year" binding:"required"`
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

	item, err := c.svc.Create(ctx.Request.Context(), teacherID, req.Degree, req.Major, req.University, req.GraduationYear)
	if err != nil {
		c.handleServiceError(ctx, log, err, "teacher-education-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateTeacherEducationController(ctx *gin.Context) {
	c.Create(ctx)
}

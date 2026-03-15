package teachersubjects

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	TeacherID string `json:"teacher_id" binding:"required"`
	SubjectID string `json:"subject_id" binding:"required"`
	Role      string `json:"role" binding:"required"`
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

	teacherID, err := uuid.Parse(req.TeacherID)
	if err != nil {
		base.BadRequest(ctx, "invalid-teacher-id", nil)
		return
	}
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		base.BadRequest(ctx, "invalid-subject-id", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), teacherID, subjectID, req.Role, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "teacher-subject-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateTeacherSubjectController(ctx *gin.Context) {
	c.Create(ctx)
}

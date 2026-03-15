package subjectgroups

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	SchoolID      string  `json:"school_id" binding:"required"`
	Code          string  `json:"code" binding:"required"`
	NameTH        string  `json:"name_th" binding:"required"`
	NameEN        *string `json:"name_en"`
	HeadTeacherID string  `json:"head_teacher_id"`
	IsActive      bool    `json:"is_active"`
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

	var headTeacherID *uuid.UUID
	if req.HeadTeacherID != "" {
		parsed, err := uuid.Parse(req.HeadTeacherID)
		if err != nil {
			base.BadRequest(ctx, "invalid-head-teacher-id", nil)
			return
		}
		headTeacherID = &parsed
	}

	item, err := c.svc.Create(ctx.Request.Context(), schoolID, req.Code, req.NameTH, req.NameEN, headTeacherID, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "subject-group-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateSubjectGroupController(ctx *gin.Context) {
	c.Create(ctx)
}

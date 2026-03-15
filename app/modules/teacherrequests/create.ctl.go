package teacherrequests

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	TeacherID       string         `json:"teacher_id" binding:"required"`
	RequestedBy     string         `json:"requested_by" binding:"required"`
	RequestedByRole string         `json:"requested_by_role" binding:"required"`
	RequestType     string         `json:"request_type" binding:"required"`
	RequestData     map[string]any `json:"request_data" binding:"required"`
	RequestReason   *string        `json:"request_reason"`
	Status          *string        `json:"status"`
	ApprovedBy      *string        `json:"approved_by"`
	ApprovedAt      *string        `json:"approved_at"`
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

	var approvedBy *uuid.UUID
	if req.ApprovedBy != nil {
		parsed, err := uuid.Parse(*req.ApprovedBy)
		if err != nil {
			base.BadRequest(ctx, "invalid-approved-by", nil)
			return
		}
		approvedBy = &parsed
	}

	item, err := c.svc.Create(ctx.Request.Context(), teacherID, req.RequestedBy, req.RequestedByRole, req.RequestType, req.RequestData, req.RequestReason, req.Status, approvedBy, req.ApprovedAt)
	if err != nil {
		c.handleServiceError(ctx, log, err, "teacher-request-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateTeacherRequestController(ctx *gin.Context) {
	c.Create(ctx)
}

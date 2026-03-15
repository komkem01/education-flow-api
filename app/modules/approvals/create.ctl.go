package approvals

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	RequestType     string         `json:"request_type" binding:"required"`
	SubjectType     string         `json:"subject_type" binding:"required"`
	SubjectID       *string        `json:"subject_id"`
	RequestedBy     string         `json:"requested_by" binding:"required"`
	RequestedByRole string         `json:"requested_by_role" binding:"required"`
	Payload         map[string]any `json:"payload" binding:"required"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), req.RequestType, req.SubjectType, req.SubjectID, req.RequestedBy, req.RequestedByRole, req.Payload)
	if err != nil {
		c.handleServiceError(ctx, log, err, "approval-request-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateApprovalRequestController(ctx *gin.Context) {
	c.Create(ctx)
}

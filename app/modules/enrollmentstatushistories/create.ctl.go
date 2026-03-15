package enrollmentstatushistories

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	EnrollmentID string  `json:"enrollment_id" binding:"required"`
	FromStatus   *string `json:"from_status"`
	ToStatus     string  `json:"to_status" binding:"required"`
	ChangedAt    *string `json:"changed_at"`
	ChangedBy    *string `json:"changed_by"`
	Reason       *string `json:"reason"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	enrollmentID, err := uuid.Parse(req.EnrollmentID)
	if err != nil {
		base.BadRequest(ctx, "invalid-enrollment-id", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), enrollmentID, req.FromStatus, req.ToStatus, req.ChangedAt, req.ChangedBy, req.Reason)
	if err != nil {
		c.handleServiceError(ctx, log, err, "enrollment-status-history-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateEnrollmentStatusHistoryController(ctx *gin.Context) {
	c.Create(ctx)
}

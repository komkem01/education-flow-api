package approvals

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ActionByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type ActionRequest struct {
	ActorID        string         `json:"actor_id" binding:"required"`
	ActorRole      string         `json:"actor_role" binding:"required"`
	IdempotencyKey *string        `json:"idempotency_key"`
	Comment        *string        `json:"comment"`
	Metadata       map[string]any `json:"metadata"`
}

func (c *Controller) Submit(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, req, ok := c.bindActionRequest(ctx)
	if !ok {
		return
	}

	item, err := c.svc.Submit(ctx.Request.Context(), id, req.ActorID, req.ActorRole, req.Comment)
	if err != nil {
		c.handleServiceError(ctx, log, err, "approval-request-submit-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) Approve(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, req, ok := c.bindActionRequest(ctx)
	if !ok {
		return
	}

	item, err := c.svc.Approve(ctx.Request.Context(), id, req.ActorID, req.ActorRole, req.IdempotencyKey, req.Comment, req.Metadata)
	if err != nil {
		c.handleServiceError(ctx, log, err, "approval-request-approve-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) Reject(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, req, ok := c.bindActionRequest(ctx)
	if !ok {
		return
	}

	item, err := c.svc.Reject(ctx.Request.Context(), id, req.ActorID, req.ActorRole, req.IdempotencyKey, req.Comment, req.Metadata)
	if err != nil {
		c.handleServiceError(ctx, log, err, "approval-request-reject-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) Cancel(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, req, ok := c.bindActionRequest(ctx)
	if !ok {
		return
	}

	item, err := c.svc.Cancel(ctx.Request.Context(), id, req.ActorID, req.ActorRole, req.IdempotencyKey, req.Comment, req.Metadata)
	if err != nil {
		c.handleServiceError(ctx, log, err, "approval-request-cancel-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) bindActionRequest(ctx *gin.Context) (uuid.UUID, *ActionRequest, bool) {
	var uriReq ActionByIDRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return uuid.Nil, nil, false
	}

	id, err := uuid.Parse(uriReq.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return uuid.Nil, nil, false
	}

	var req ActionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return uuid.Nil, nil, false
	}

	return id, &req, true
}

func (c *Controller) ApprovalRequestsSubmit(ctx *gin.Context) {
	c.Submit(ctx)
}

func (c *Controller) ApprovalRequestsApprove(ctx *gin.Context) {
	c.Approve(ctx)
}

func (c *Controller) ApprovalRequestsReject(ctx *gin.Context) {
	c.Reject(ctx)
}

func (c *Controller) ApprovalRequestsCancel(ctx *gin.Context) {
	c.Cancel(ctx)
}

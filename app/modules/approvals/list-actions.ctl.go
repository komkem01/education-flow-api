package approvals

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListActionsByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type ListActionsRequest struct {
	base.RequestPaginate
	ActedBy string `form:"acted_by"`
	Action  string `form:"action"`
}

func (c *Controller) ListActions(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var uriReq ListActionsByIDRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	requestID, err := uuid.Parse(uriReq.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return
	}

	var req ListActionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var actedBy *uuid.UUID
	if req.ActedBy != "" {
		parsed, err := uuid.Parse(req.ActedBy)
		if err != nil {
			base.BadRequest(ctx, "invalid-acted-by", nil)
			return
		}
		actedBy = &parsed
	}

	var action *string
	if req.Action != "" {
		action = &req.Action
	}

	items, page, err := c.svc.ListActions(ctx.Request.Context(), &req.RequestPaginate, requestID, actedBy, action)
	if err != nil {
		c.handleServiceError(ctx, log, err, "approval-actions-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) ApprovalRequestsActionsList(ctx *gin.Context) {
	c.ListActions(ctx)
}

package approvals

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	RequestedBy     string `form:"requested_by"`
	RequestedByRole string `form:"requested_by_role"`
	Status          string `form:"status"`
	RequestType     string `form:"request_type"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var requestedBy *uuid.UUID
	if req.RequestedBy != "" {
		parsed, err := uuid.Parse(req.RequestedBy)
		if err != nil {
			base.BadRequest(ctx, "invalid-requested-by", nil)
			return
		}
		requestedBy = &parsed
	}

	var requestedByRole *string
	if req.RequestedByRole != "" {
		requestedByRole = &req.RequestedByRole
	}

	var status *string
	if req.Status != "" {
		status = &req.Status
	}

	var requestType *string
	if req.RequestType != "" {
		requestType = &req.RequestType
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, requestedBy, requestedByRole, status, requestType)
	if err != nil {
		c.handleServiceError(ctx, log, err, "approval-request-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) ApprovalRequestsList(ctx *gin.Context) {
	c.List(ctx)
}

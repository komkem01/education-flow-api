package enrollmentstatushistories

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	EnrollmentID string `form:"enrollment_id"`
	ToStatus     string `form:"to_status"`
	ChangedBy    string `form:"changed_by"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	enrollmentID, err := parseUUIDFilter(req.EnrollmentID)
	if err != nil {
		base.BadRequest(ctx, "invalid-enrollment-id", nil)
		return
	}
	changedBy, err := parseUUIDFilter(req.ChangedBy)
	if err != nil {
		base.BadRequest(ctx, "invalid-changed-by", nil)
		return
	}

	var toStatus *string
	if req.ToStatus != "" {
		toStatus = &req.ToStatus
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, enrollmentID, toStatus, changedBy)
	if err != nil {
		c.handleServiceError(ctx, log, err, "enrollment-status-history-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) EnrollmentStatusHistoriesList(ctx *gin.Context) {
	c.List(ctx)
}

func parseUUIDFilter(v string) (*uuid.UUID, error) {
	if v == "" {
		return nil, nil
	}
	parsed, err := uuid.Parse(v)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

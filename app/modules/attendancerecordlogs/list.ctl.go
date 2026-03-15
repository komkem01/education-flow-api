package attendancerecordlogs

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	RecordID  string `form:"record_id"`
	ChangedBy string `form:"changed_by"`
	NewStatus string `form:"new_status"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	recordID, err := parseUUIDFilter(req.RecordID)
	if err != nil {
		base.BadRequest(ctx, "invalid-record-id", nil)
		return
	}
	changedBy, err := parseUUIDFilter(req.ChangedBy)
	if err != nil {
		base.BadRequest(ctx, "invalid-changed-by", nil)
		return
	}

	var newStatus *string
	if req.NewStatus != "" {
		newStatus = &req.NewStatus
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, recordID, changedBy, newStatus)
	if err != nil {
		c.handleServiceError(ctx, log, err, "attendance-record-log-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) AttendanceRecordLogsList(ctx *gin.Context) {
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

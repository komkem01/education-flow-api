package attendancerecords

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	SessionID    string `form:"session_id"`
	EnrollmentID string `form:"enrollment_id"`
	Status       string `form:"status"`
	Source       string `form:"source"`
	MarkedBy     string `form:"marked_by"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	sessionID, err := parseUUIDFilter(req.SessionID)
	if err != nil {
		base.BadRequest(ctx, "invalid-session-id", nil)
		return
	}
	enrollmentID, err := parseUUIDFilter(req.EnrollmentID)
	if err != nil {
		base.BadRequest(ctx, "invalid-enrollment-id", nil)
		return
	}
	markedBy, err := parseUUIDFilter(req.MarkedBy)
	if err != nil {
		base.BadRequest(ctx, "invalid-marked-by", nil)
		return
	}

	var status *string
	if req.Status != "" {
		status = &req.Status
	}
	var source *string
	if req.Source != "" {
		source = &req.Source
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, sessionID, enrollmentID, status, source, markedBy)
	if err != nil {
		c.handleServiceError(ctx, log, err, "attendance-record-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) AttendanceRecordsList(ctx *gin.Context) {
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

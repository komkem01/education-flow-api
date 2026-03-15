package auditlogs

import (
	"strings"
	"time"

	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	ActorID    string `form:"actor_id"`
	ActorRole  string `form:"actor_role"`
	Method     string `form:"method"`
	Path       string `form:"path"`
	StatusCode *int   `form:"status_code"`
	From       string `form:"from"`
	To         string `form:"to"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	var actorID *uuid.UUID
	if req.ActorID != "" {
		parsed, err := uuid.Parse(req.ActorID)
		if err != nil {
			base.BadRequest(ctx, "invalid-request-form", nil)
			return
		}
		actorID = &parsed
	}

	var actorRole *string
	if req.ActorRole != "" {
		v := strings.TrimSpace(req.ActorRole)
		actorRole = &v
	}

	var method *string
	if req.Method != "" {
		v := strings.TrimSpace(req.Method)
		method = &v
	}

	var path *string
	if req.Path != "" {
		v := strings.TrimSpace(req.Path)
		path = &v
	}

	from, err := parseQueryTime(req.From)
	if err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	to, err := parseQueryTime(req.To)
	if err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	if from != nil && to != nil && from.After(*to) {
		base.BadRequest(ctx, "invalid-date-range", nil)
		return
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, actorID, actorRole, method, path, req.StatusCode, from, to)
	if err != nil {
		c.handleServiceError(ctx, log, err, "audit-log-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) AuditLogsList(ctx *gin.Context) {
	c.List(ctx)
}

func parseQueryTime(v string) (*time.Time, error) {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, trimmed)
	if err == nil {
		return &parsed, nil
	}
	parsed, err = time.Parse("2006-01-02", trimmed)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

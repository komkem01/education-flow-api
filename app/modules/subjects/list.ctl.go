package subjects

import (
	"strconv"

	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	IsActive       *bool  `form:"is_active"`
	SchoolID       string `form:"school_id"`
	SubjectGroupID string `form:"subject_group_id"`
	IsElective     *bool  `form:"is_elective"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	if raw := ctx.Query("is_active"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			base.BadRequest(ctx, "invalid-is-active", nil)
			return
		}
		req.IsActive = &parsed
	}
	if raw := ctx.Query("is_elective"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			base.BadRequest(ctx, "invalid-is-elective", nil)
			return
		}
		req.IsElective = &parsed
	}

	var schoolID *uuid.UUID
	if req.SchoolID != "" {
		parsed, err := uuid.Parse(req.SchoolID)
		if err != nil {
			base.BadRequest(ctx, "invalid-school-id", nil)
			return
		}
		schoolID = &parsed
	}

	var subjectGroupID *uuid.UUID
	if req.SubjectGroupID != "" {
		parsed, err := uuid.Parse(req.SubjectGroupID)
		if err != nil {
			base.BadRequest(ctx, "invalid-subject-group-id", nil)
			return
		}
		subjectGroupID = &parsed
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, req.IsActive, schoolID, subjectGroupID, req.IsElective)
	if err != nil {
		c.handleServiceError(ctx, log, err, "subject-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) SubjectsList(ctx *gin.Context) {
	c.List(ctx)
}

package classrooms

import (
	"strconv"

	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	IsActive          *bool  `form:"is_active"`
	SchoolID          string `form:"school_id"`
	AcademicYearID    string `form:"academic_year_id"`
	HomeroomTeacherID string `form:"homeroom_teacher_id"`
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

	var schoolID *uuid.UUID
	if req.SchoolID != "" {
		parsed, err := uuid.Parse(req.SchoolID)
		if err != nil {
			base.BadRequest(ctx, "invalid-school-id", nil)
			return
		}
		schoolID = &parsed
	}

	var academicYearID *uuid.UUID
	if req.AcademicYearID != "" {
		parsed, err := uuid.Parse(req.AcademicYearID)
		if err != nil {
			base.BadRequest(ctx, "invalid-academic-year-id", nil)
			return
		}
		academicYearID = &parsed
	}

	var homeroomTeacherID *uuid.UUID
	if req.HomeroomTeacherID != "" {
		parsed, err := uuid.Parse(req.HomeroomTeacherID)
		if err != nil {
			base.BadRequest(ctx, "invalid-homeroom-teacher-id", nil)
			return
		}
		homeroomTeacherID = &parsed
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, req.IsActive, schoolID, academicYearID, homeroomTeacherID)
	if err != nil {
		c.handleServiceError(ctx, log, err, "classroom-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) ClassroomsList(ctx *gin.Context) {
	c.List(ctx)
}

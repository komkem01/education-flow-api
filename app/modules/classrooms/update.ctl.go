package classrooms

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateRequest struct {
	SchoolID          *string `json:"school_id"`
	AcademicYearID    *string `json:"academic_year_id"`
	Level             *string `json:"level"`
	GradeLevel        *string `json:"grade_level"`
	RoomNo            *string `json:"room_no"`
	Name              *string `json:"name"`
	HomeroomTeacherID *string `json:"homeroom_teacher_id"`
	AdvisorTeacherID  *string `json:"advisor_teacher_id"`
	Capacity          *int    `json:"capacity"`
	IsActive          *bool   `json:"is_active"`
}

func (c *Controller) Update(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var uriReq UpdateByIDRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	id, err := uuid.Parse(uriReq.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return
	}

	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	if req.Level == nil && req.GradeLevel != nil {
		req.Level = req.GradeLevel
	}
	if req.HomeroomTeacherID == nil && req.AdvisorTeacherID != nil {
		req.HomeroomTeacherID = req.AdvisorTeacherID
	}

	item, err := c.svc.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		c.handleServiceError(ctx, log, err, "classroom-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) ClassroomsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

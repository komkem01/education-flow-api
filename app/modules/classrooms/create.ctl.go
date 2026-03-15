package classrooms

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	SchoolID          string `json:"school_id" binding:"required"`
	AcademicYearID    string `json:"academic_year_id" binding:"required"`
	Level             string `json:"level" binding:"required"`
	RoomNo            string `json:"room_no" binding:"required"`
	Name              string `json:"name" binding:"required"`
	HomeroomTeacherID string `json:"homeroom_teacher_id"`
	Capacity          *int   `json:"capacity"`
	IsActive          bool   `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, "invalid-school-id", nil)
		return
	}
	academicYearID, err := uuid.Parse(req.AcademicYearID)
	if err != nil {
		base.BadRequest(ctx, "invalid-academic-year-id", nil)
		return
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

	item, err := c.svc.Create(ctx.Request.Context(), schoolID, academicYearID, req.Level, req.RoomNo, req.Name, homeroomTeacherID, req.Capacity, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "classroom-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateClassroomController(ctx *gin.Context) {
	c.Create(ctx)
}

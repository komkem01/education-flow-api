package classrooms

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	SchoolID          string `json:"school_id"`
	AcademicYearID    string `json:"academic_year_id"`
	Level             string `json:"level"`
	GradeLevel        string `json:"grade_level"`
	RoomNo            string `json:"room_no"`
	Name              string `json:"name"`
	HomeroomTeacherID string `json:"homeroom_teacher_id"`
	AdvisorTeacherID  string `json:"advisor_teacher_id"`
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

	if req.Level == "" && req.GradeLevel != "" {
		req.Level = req.GradeLevel
	}
	if req.HomeroomTeacherID == "" && req.AdvisorTeacherID != "" {
		req.HomeroomTeacherID = req.AdvisorTeacherID
	}

	if req.SchoolID == "" {
		base.BadRequest(ctx, "invalid-school-id", nil)
		return
	}
	if req.AcademicYearID == "" {
		base.BadRequest(ctx, "invalid-academic-year-id", nil)
		return
	}
	if req.Level == "" {
		base.BadRequest(ctx, "invalid-level", nil)
		return
	}
	if req.RoomNo == "" {
		base.BadRequest(ctx, "invalid-room-no", nil)
		return
	}
	if req.Name == "" {
		base.BadRequest(ctx, "invalid-name", nil)
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

package memberstudents

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	MemberID         string  `json:"member_id" binding:"required"`
	SchoolID         string  `json:"school_id" binding:"required"`
	GenderID         string  `json:"gender_id" binding:"required"`
	PrefixID         string  `json:"prefix_id" binding:"required"`
	AdvisorTeacherID string  `json:"advisor_teacher_id"`
	StudentCode      string  `json:"student_code" binding:"required"`
	FirstNameTH      string  `json:"first_name_th" binding:"required"`
	LastNameTH       string  `json:"last_name_th" binding:"required"`
	FirstNameEN      *string `json:"first_name_en"`
	LastNameEN       *string `json:"last_name_en"`
	CitizenID        *string `json:"citizen_id"`
	Phone            *string `json:"phone"`
	IsActive         bool    `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	memberID, err := uuid.Parse(req.MemberID)
	if err != nil {
		base.BadRequest(ctx, "invalid-member-id", nil)
		return
	}
	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, "invalid-school-id", nil)
		return
	}
	genderID, err := uuid.Parse(req.GenderID)
	if err != nil {
		base.BadRequest(ctx, "invalid-gender-id", nil)
		return
	}
	prefixID, err := uuid.Parse(req.PrefixID)
	if err != nil {
		base.BadRequest(ctx, "invalid-prefix-id", nil)
		return
	}

	var advisorTeacherID *uuid.UUID
	if req.AdvisorTeacherID != "" {
		parsed, err := uuid.Parse(req.AdvisorTeacherID)
		if err != nil {
			base.BadRequest(ctx, "invalid-advisor-teacher-id", nil)
			return
		}
		advisorTeacherID = &parsed
	}

	item, err := c.svc.Create(ctx.Request.Context(), memberID, schoolID, genderID, prefixID, advisorTeacherID, req.StudentCode, req.FirstNameTH, req.LastNameTH, req.FirstNameEN, req.LastNameEN, req.CitizenID, req.Phone, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-student-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateMemberStudentController(ctx *gin.Context) {
	c.Create(ctx)
}

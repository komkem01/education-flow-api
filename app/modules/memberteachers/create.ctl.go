package memberteachers

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	MemberID         string  `json:"member_id" binding:"required"`
	GenderID         string  `json:"gender_id" binding:"required"`
	PrefixID         string  `json:"prefix_id" binding:"required"`
	Code             string  `json:"code" binding:"required"`
	CitizenID        string  `json:"citizen_id" binding:"required"`
	FirstNameTH      string  `json:"first_name_th" binding:"required"`
	LastNameTH       string  `json:"last_name_th" binding:"required"`
	FirstNameEN      string  `json:"first_name_en" binding:"required"`
	LastNameEN       string  `json:"last_name_en" binding:"required"`
	Phone            string  `json:"phone" binding:"required"`
	Position         string  `json:"position" binding:"required"`
	AcademicStanding string  `json:"academic_standing" binding:"required"`
	DepartmentID     string  `json:"department_id" binding:"required"`
	StartDate        string  `json:"start_date" binding:"required"`
	EndDate          *string `json:"end_date"`
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
	departmentID, err := uuid.Parse(req.DepartmentID)
	if err != nil {
		base.BadRequest(ctx, "invalid-department-id", nil)
		return
	}

	teacher, err := c.svc.Create(ctx.Request.Context(), memberID, genderID, prefixID, req.Code, req.CitizenID, req.FirstNameTH, req.LastNameTH, req.FirstNameEN, req.LastNameEN, req.Phone, req.Position, req.AcademicStanding, departmentID, req.StartDate, req.EndDate, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-teacher-create-failed")
		return
	}

	base.Success(ctx, teacher, "success")
}

func (c *Controller) CreateMemberTeacherController(ctx *gin.Context) {
	c.Create(ctx)
}

func (c *Controller) MemberTeachersCreate(ctx *gin.Context) {
	c.Create(ctx)
}

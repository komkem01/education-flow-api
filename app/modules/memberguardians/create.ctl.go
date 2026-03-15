package memberguardians

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	MemberID    string  `json:"member_id"`
	SchoolID    string  `json:"school_id" binding:"required"`
	GenderID    string  `json:"gender_id" binding:"required"`
	PrefixID    string  `json:"prefix_id" binding:"required"`
	FirstNameTH string  `json:"first_name_th" binding:"required"`
	LastNameTH  string  `json:"last_name_th" binding:"required"`
	FirstNameEN *string `json:"first_name_en"`
	LastNameEN  *string `json:"last_name_en"`
	CitizenID   *string `json:"citizen_id"`
	Phone       *string `json:"phone"`
	IsActive    bool    `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var memberID *uuid.UUID
	if req.MemberID != "" {
		parsed, err := uuid.Parse(req.MemberID)
		if err != nil {
			base.BadRequest(ctx, "invalid-member-id", nil)
			return
		}
		memberID = &parsed
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

	item, err := c.svc.Create(ctx.Request.Context(), memberID, schoolID, genderID, prefixID, req.FirstNameTH, req.LastNameTH, req.FirstNameEN, req.LastNameEN, req.CitizenID, req.Phone, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-guardian-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateMemberGuardianController(ctx *gin.Context) {
	c.Create(ctx)
}

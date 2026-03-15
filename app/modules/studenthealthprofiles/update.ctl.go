package studenthealthprofiles

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
	StudentID      *string `json:"student_id"`
	BloodType      *string `json:"blood_type"`
	AllergyInfo    *string `json:"allergy_info"`
	ChronicDisease *string `json:"chronic_disease"`
	MedicalNote    *string `json:"medical_note"`
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

	item, err := c.svc.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-health-profile-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) StudentHealthProfilesUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

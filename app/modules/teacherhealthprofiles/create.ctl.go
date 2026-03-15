package teacherhealthprofiles

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	MemberTeacherID    string  `json:"member_teacher_id" binding:"required"`
	BloodType          *string `json:"blood_type"`
	AllergyInfo        *string `json:"allergy_info"`
	ChronicDisease     *string `json:"chronic_disease"`
	MedicationNote     *string `json:"medication_note"`
	FitnessForWorkNote *string `json:"fitness_for_work_note"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	memberTeacherID, err := uuid.Parse(req.MemberTeacherID)
	if err != nil {
		base.BadRequest(ctx, "invalid-member-teacher-id", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), memberTeacherID, req.BloodType, req.AllergyInfo, req.ChronicDisease, req.MedicationNote, req.FitnessForWorkNote)
	if err != nil {
		c.handleServiceError(ctx, log, err, "teacher-health-profile-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateTeacherHealthProfileController(ctx *gin.Context) {
	c.Create(ctx)
}

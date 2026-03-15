package studenthealthprofiles

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	StudentID      string  `json:"student_id" binding:"required"`
	BloodType      *string `json:"blood_type"`
	AllergyInfo    *string `json:"allergy_info"`
	ChronicDisease *string `json:"chronic_disease"`
	MedicalNote    *string `json:"medical_note"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		base.BadRequest(ctx, "invalid-student-id", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), studentID, req.BloodType, req.AllergyInfo, req.ChronicDisease, req.MedicalNote)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-health-profile-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateStudentHealthProfileController(ctx *gin.Context) {
	c.Create(ctx)
}

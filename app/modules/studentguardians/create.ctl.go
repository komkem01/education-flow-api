package studentguardians

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	StudentID          string  `json:"student_id" binding:"required"`
	GuardianID         string  `json:"guardian_id" binding:"required"`
	Relationship       string  `json:"relationship" binding:"required"`
	IsMainGuardian     bool    `json:"is_main_guardian"`
	CanPickup          bool    `json:"can_pickup"`
	IsEmergencyContact bool    `json:"is_emergency_contact"`
	Note               *string `json:"note"`
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
	guardianID, err := uuid.Parse(req.GuardianID)
	if err != nil {
		base.BadRequest(ctx, "invalid-guardian-id", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), studentID, guardianID, req.Relationship, req.IsMainGuardian, req.CanPickup, req.IsEmergencyContact, req.Note)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-guardian-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateStudentGuardianController(ctx *gin.Context) {
	c.Create(ctx)
}

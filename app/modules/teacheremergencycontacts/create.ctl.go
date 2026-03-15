package teacheremergencycontacts

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	MemberTeacherID      string  `json:"member_teacher_id" binding:"required"`
	EmergencyContactName string  `json:"emergency_contact_name" binding:"required"`
	Relationship         string  `json:"relationship" binding:"required"`
	PhonePrimary         string  `json:"phone_primary" binding:"required"`
	PhoneSecondary       *string `json:"phone_secondary"`
	CanDecideMedical     bool    `json:"can_decide_medical"`
	IsPrimary            bool    `json:"is_primary"`
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

	item, err := c.svc.Create(ctx.Request.Context(), memberTeacherID, req.EmergencyContactName, req.Relationship, req.PhonePrimary, req.PhoneSecondary, req.CanDecideMedical, req.IsPrimary)
	if err != nil {
		c.handleServiceError(ctx, log, err, "teacher-emergency-contact-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateTeacherEmergencyContactController(ctx *gin.Context) {
	c.Create(ctx)
}

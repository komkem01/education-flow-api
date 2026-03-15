package teacheremergencycontacts

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
	EmergencyContactName *string `json:"emergency_contact_name"`
	Relationship         *string `json:"relationship"`
	PhonePrimary         *string `json:"phone_primary"`
	PhoneSecondary       *string `json:"phone_secondary"`
	CanDecideMedical     *bool   `json:"can_decide_medical"`
	IsPrimary            *bool   `json:"is_primary"`
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
		c.handleServiceError(ctx, log, err, "teacher-emergency-contact-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) TeacherEmergencyContactsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

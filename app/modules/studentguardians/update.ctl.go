package studentguardians

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
	StudentID          *string `json:"student_id"`
	GuardianID         *string `json:"guardian_id"`
	Relationship       *string `json:"relationship"`
	IsMainGuardian     *bool   `json:"is_main_guardian"`
	CanPickup          *bool   `json:"can_pickup"`
	IsEmergencyContact *bool   `json:"is_emergency_contact"`
	Note               *string `json:"note"`
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
		c.handleServiceError(ctx, log, err, "student-guardian-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) StudentGuardiansUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

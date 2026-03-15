package teacherlicenses

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
	TeacherID     *string `json:"teacher_id"`
	LicenseNo     *string `json:"license_no"`
	IssuedAt      *string `json:"issued_at"`
	ExpiresAt     *string `json:"expires_at"`
	LicenseStatus *string `json:"license_status"`
	IssuedBy      *string `json:"issued_by"`
	Note          *string `json:"note"`
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
		c.handleServiceError(ctx, log, err, "teacher-license-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) TeacherLicensesUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

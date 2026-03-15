package teacherlicenses

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	TeacherID     string  `json:"teacher_id" binding:"required"`
	LicenseNo     string  `json:"license_no" binding:"required"`
	IssuedAt      *string `json:"issued_at"`
	ExpiresAt     *string `json:"expires_at"`
	LicenseStatus string  `json:"license_status" binding:"required"`
	IssuedBy      *string `json:"issued_by"`
	Note          *string `json:"note"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	teacherID, err := uuid.Parse(req.TeacherID)
	if err != nil {
		base.BadRequest(ctx, "invalid-teacher-id", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), teacherID, req.LicenseNo, req.IssuedAt, req.ExpiresAt, req.LicenseStatus, req.IssuedBy, req.Note)
	if err != nil {
		c.handleServiceError(ctx, log, err, "teacher-license-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateTeacherLicenseController(ctx *gin.Context) {
	c.Create(ctx)
}

package studentprofiles

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	StudentID             string  `json:"student_id" binding:"required"`
	BirthDate             *string `json:"birth_date"`
	Nationality           *string `json:"nationality"`
	Religion              *string `json:"religion"`
	AddressCurrent        *string `json:"address_current"`
	AddressRegistered     *string `json:"address_registered"`
	EmergencyContactName  *string `json:"emergency_contact_name"`
	EmergencyContactPhone *string `json:"emergency_contact_phone"`
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

	item, err := c.svc.Create(ctx.Request.Context(), studentID, req.BirthDate, req.Nationality, req.Religion, req.AddressCurrent, req.AddressRegistered, req.EmergencyContactName, req.EmergencyContactPhone)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-profile-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateStudentProfileController(ctx *gin.Context) {
	c.Create(ctx)
}

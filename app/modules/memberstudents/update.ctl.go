package memberstudents

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
	MemberID         *string `json:"member_id"`
	SchoolID         *string `json:"school_id"`
	GenderID         *string `json:"gender_id"`
	PrefixID         *string `json:"prefix_id"`
	AdvisorTeacherID *string `json:"advisor_teacher_id"`
	StudentCode      *string `json:"student_code"`
	FirstNameTH      *string `json:"first_name_th"`
	LastNameTH       *string `json:"last_name_th"`
	FirstNameEN      *string `json:"first_name_en"`
	LastNameEN       *string `json:"last_name_en"`
	CitizenID        *string `json:"citizen_id"`
	Phone            *string `json:"phone"`
	IsActive         *bool   `json:"is_active"`
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
		c.handleServiceError(ctx, log, err, "member-student-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) MemberStudentsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

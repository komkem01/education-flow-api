package studentenrollments

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
	StudentID            *string `json:"student_id"`
	SchoolID             *string `json:"school_id"`
	AcademicYearID       *string `json:"academic_year_id"`
	ClassroomID          *string `json:"classroom_id"`
	EnrolledAt           *string `json:"enrolled_at"`
	ExitedAt             *string `json:"exited_at"`
	Status               *string `json:"status"`
	EnrollmentType       *string `json:"enrollment_type"`
	ExitReason           *string `json:"exit_reason"`
	ExitNote             *string `json:"exit_note"`
	PreviousEnrollmentID *string `json:"previous_enrollment_id"`
	RollNo               *string `json:"roll_no"`
	ApprovedBy           *string `json:"approved_by"`
	ApprovedAt           *string `json:"approved_at"`
	ApprovalNote         *string `json:"approval_note"`
	CreatedBy            *string `json:"created_by"`
	UpdatedBy            *string `json:"updated_by"`
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
		c.handleServiceError(ctx, log, err, "student-enrollment-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) StudentEnrollmentsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

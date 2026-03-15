package studentenrollments

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	StudentID            string  `json:"student_id" binding:"required"`
	SchoolID             string  `json:"school_id" binding:"required"`
	AcademicYearID       string  `json:"academic_year_id" binding:"required"`
	ClassroomID          string  `json:"classroom_id" binding:"required"`
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
	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, "invalid-school-id", nil)
		return
	}
	academicYearID, err := uuid.Parse(req.AcademicYearID)
	if err != nil {
		base.BadRequest(ctx, "invalid-academic-year-id", nil)
		return
	}
	classroomID, err := uuid.Parse(req.ClassroomID)
	if err != nil {
		base.BadRequest(ctx, "invalid-classroom-id", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), studentID, schoolID, academicYearID, classroomID, req.EnrolledAt, req.ExitedAt, req.Status, req.EnrollmentType, req.ExitReason, req.ExitNote, req.PreviousEnrollmentID, req.RollNo, req.ApprovedBy, req.ApprovedAt, req.ApprovalNote, req.CreatedBy, req.UpdatedBy)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-enrollment-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateStudentEnrollmentController(ctx *gin.Context) {
	c.Create(ctx)
}

package enrollmentsubjects

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	EnrollmentID  string   `json:"enrollment_id" binding:"required"`
	SubjectID     string   `json:"subject_id" binding:"required"`
	TeacherID     *string  `json:"teacher_id"`
	IsPrimary     bool     `json:"is_primary"`
	Status        *string  `json:"status"`
	MidtermScore  *float64 `json:"midterm_score"`
	FinalScore    *float64 `json:"final_score"`
	ActivityScore *float64 `json:"activity_score"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	enrollmentID, err := uuid.Parse(req.EnrollmentID)
	if err != nil {
		base.BadRequest(ctx, "invalid-enrollment-id", nil)
		return
	}
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		base.BadRequest(ctx, "invalid-subject-id", nil)
		return
	}

	item, err := c.svc.Create(
		ctx.Request.Context(),
		enrollmentID,
		subjectID,
		req.TeacherID,
		req.IsPrimary,
		req.Status,
		req.MidtermScore,
		req.FinalScore,
		req.ActivityScore,
	)
	if err != nil {
		c.handleServiceError(ctx, log, err, "enrollment-subject-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateEnrollmentSubjectController(ctx *gin.Context) {
	c.Create(ctx)
}

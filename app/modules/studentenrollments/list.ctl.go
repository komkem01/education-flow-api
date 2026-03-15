package studentenrollments

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	StudentID      string `form:"student_id"`
	SchoolID       string `form:"school_id"`
	AcademicYearID string `form:"academic_year_id"`
	ClassroomID    string `form:"classroom_id"`
	Status         string `form:"status"`
	EnrollmentType string `form:"enrollment_type"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	studentID, err := parseUUIDFilter(req.StudentID)
	if err != nil {
		base.BadRequest(ctx, "invalid-student-id", nil)
		return
	}
	schoolID, err := parseUUIDFilter(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, "invalid-school-id", nil)
		return
	}
	academicYearID, err := parseUUIDFilter(req.AcademicYearID)
	if err != nil {
		base.BadRequest(ctx, "invalid-academic-year-id", nil)
		return
	}
	classroomID, err := parseUUIDFilter(req.ClassroomID)
	if err != nil {
		base.BadRequest(ctx, "invalid-classroom-id", nil)
		return
	}

	var status *string
	if req.Status != "" {
		status = &req.Status
	}
	var enrollmentType *string
	if req.EnrollmentType != "" {
		enrollmentType = &req.EnrollmentType
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, studentID, schoolID, academicYearID, classroomID, status, enrollmentType)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-enrollment-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) StudentEnrollmentsList(ctx *gin.Context) {
	c.List(ctx)
}

func parseUUIDFilter(v string) (*uuid.UUID, error) {
	if v == "" {
		return nil, nil
	}
	parsed, err := uuid.Parse(v)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

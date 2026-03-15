package enrollmentsubjects

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	EnrollmentID string `form:"enrollment_id"`
	SubjectID    string `form:"subject_id"`
	TeacherID    string `form:"teacher_id"`
	Status       string `form:"status"`
	IsPrimary    string `form:"is_primary"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	enrollmentID, err := parseUUIDFilter(req.EnrollmentID)
	if err != nil {
		base.BadRequest(ctx, "invalid-enrollment-id", nil)
		return
	}
	subjectID, err := parseUUIDFilter(req.SubjectID)
	if err != nil {
		base.BadRequest(ctx, "invalid-subject-id", nil)
		return
	}
	teacherID, err := parseUUIDFilter(req.TeacherID)
	if err != nil {
		base.BadRequest(ctx, "invalid-teacher-id", nil)
		return
	}

	var status *string
	if req.Status != "" {
		status = &req.Status
	}

	var isPrimary *bool
	if req.IsPrimary != "" {
		parsed, ok := parseBoolFilter(req.IsPrimary)
		if !ok {
			base.BadRequest(ctx, "invalid-is-primary", nil)
			return
		}
		isPrimary = &parsed
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, enrollmentID, subjectID, teacherID, status, isPrimary)
	if err != nil {
		c.handleServiceError(ctx, log, err, "enrollment-subject-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) EnrollmentSubjectsList(ctx *gin.Context) {
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

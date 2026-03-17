package reports

import (
	"strconv"

	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (c *Controller) Summary(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	user, ok := c.currentUser(ctx)
	if !ok {
		return
	}

	if !c.canRead(user) {
		base.Forbidden(ctx, "forbidden", nil)
		return
	}

	var academicYearID *uuid.UUID
	if raw := ctx.Query("academic_year_id"); raw != "" {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			base.BadRequest(ctx, "invalid-academic-year-id", nil)
			return
		}
		academicYearID = &parsed
	}

	var semesterNo *int
	if raw := ctx.Query("semester_no"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			base.BadRequest(ctx, "invalid-semester-no", nil)
			return
		}
		semesterNo = &parsed
	}

	res, err := c.svc.Summary(ctx.Request.Context(), fallbackSchoolID(user.Member), &SummaryRequest{
		AcademicYearID: academicYearID,
		SemesterNo:     semesterNo,
	})
	if err != nil {
		log.Errf("report-summary-failed: %v", err)
		base.InternalServerError(ctx, "report-summary-failed", nil)
		return
	}

	base.Success(ctx, res, "success")
}

func (c *Controller) ReportsSummary(ctx *gin.Context) {
	c.Summary(ctx)
}

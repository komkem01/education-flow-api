package reports

import (
	"strconv"

	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (c *Controller) Info(ctx *gin.Context) {
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

	res, err := c.svc.Info(ctx.Request.Context(), fallbackSchoolID(user.Member), ctx.Param("id"), &SummaryRequest{
		AcademicYearID: academicYearID,
		SemesterNo:     semesterNo,
	})
	if err != nil {
		log.Errf("report-info-failed: %v", err)
		base.InternalServerError(ctx, "report-info-failed", nil)
		return
	}

	base.Success(ctx, res, "success")
}

func (c *Controller) ReportsInfo(ctx *gin.Context) {
	c.Info(ctx)
}

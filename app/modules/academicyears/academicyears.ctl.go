package academicyears

import (
	"errors"

	"eduflow/app/utils/base"
	logpkg "eduflow/internal/log"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{
		tracer: trace,
		svc:    svc,
	}
}

func (c *Controller) handleServiceError(ctx *gin.Context, log *logpkg.Logger, err error, fallback string) {
	switch {
	case errors.Is(err, ErrAcademicYearNotFound):
		base.BadRequest(ctx, "academic-year-not-found", nil)
	case errors.Is(err, ErrAcademicYearUnauthorized):
		base.Unauthorized(ctx, "unauthorized", nil)
	case errors.Is(err, ErrAcademicYearConditionFail):
		base.ValidateFailed(ctx, "condition-fail", nil)
	case errors.Is(err, ErrAcademicYearDuplicate):
		base.BadRequest(ctx, "academic-year-duplicate", nil)
	case errors.Is(err, ErrAcademicYearInvalidUpdate):
		base.ValidateFailed(ctx, "condition-fail", nil)
	default:
		log.Errf("%s: %v", fallback, err)
		base.InternalServerError(ctx, fallback, nil)
	}
}

package memberteachers

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
	return &Controller{tracer: trace, svc: svc}
}

func (c *Controller) handleServiceError(ctx *gin.Context, log *logpkg.Logger, err error, fallback string) {
	switch {
	case errors.Is(err, ErrMemberTeacherNotFound):
		base.BadRequest(ctx, "member-teacher-not-found", nil)
	case errors.Is(err, ErrMemberTeacherUnauthorized):
		base.Unauthorized(ctx, "unauthorized", nil)
	case errors.Is(err, ErrTeacherInvalidEmail):
		base.ValidateFailed(ctx, "invalid-email", nil)
	case errors.Is(err, ErrTeacherInvalidPassword):
		base.ValidateFailed(ctx, "invalid-password", nil)
	case errors.Is(err, ErrTeacherInvalidCitizenID):
		base.ValidateFailed(ctx, "invalid-citizen-id", nil)
	case errors.Is(err, ErrTeacherInvalidPhone):
		base.ValidateFailed(ctx, "invalid-phone", nil)
	case errors.Is(err, ErrTeacherInvalidDateRange):
		base.ValidateFailed(ctx, "invalid-date-range", nil)
	case errors.Is(err, ErrTeacherAddressPrimaryDup):
		base.ValidateFailed(ctx, "teacher-address-primary-duplicate", nil)
	case errors.Is(err, ErrMemberTeacherConditionFail):
		base.ValidateFailed(ctx, "condition-fail", nil)
	case errors.Is(err, ErrMemberTeacherDuplicate):
		base.BadRequest(ctx, "member-teacher-duplicate", nil)
	default:
		log.Errf("%s: %v", fallback, err)
		base.InternalServerError(ctx, fallback, nil)
	}
}

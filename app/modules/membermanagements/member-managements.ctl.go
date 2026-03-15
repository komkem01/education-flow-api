package membermanagements

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
	case errors.Is(err, ErrMemberManagementNotFound):
		base.BadRequest(ctx, "member-management-not-found", nil)
	case errors.Is(err, ErrMemberManagementUnauthorized):
		base.Unauthorized(ctx, "unauthorized", nil)
	case errors.Is(err, ErrManagementInvalidEmail):
		base.ValidateFailed(ctx, "invalid-email", nil)
	case errors.Is(err, ErrManagementInvalidPassword):
		base.ValidateFailed(ctx, "invalid-password", nil)
	case errors.Is(err, ErrManagementInvalidReason):
		base.ValidateFailed(ctx, "invalid-request-reason", nil)
	case errors.Is(err, ErrMemberManagementConditionFail):
		base.ValidateFailed(ctx, "condition-fail", nil)
	case errors.Is(err, ErrMemberManagementDuplicate):
		base.BadRequest(ctx, "member-management-duplicate", nil)
	default:
		log.Errf("%s: %v", fallback, err)
		base.InternalServerError(ctx, fallback, nil)
	}
}

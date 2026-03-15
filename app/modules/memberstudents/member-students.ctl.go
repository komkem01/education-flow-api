package memberstudents

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
	case errors.Is(err, ErrMemberStudentNotFound):
		base.BadRequest(ctx, "member-student-not-found", nil)
	case errors.Is(err, ErrMemberStudentUnauthorized):
		base.Unauthorized(ctx, "unauthorized", nil)
	case errors.Is(err, ErrInvalidEmail):
		base.ValidateFailed(ctx, "invalid-email", nil)
	case errors.Is(err, ErrInvalidPassword):
		base.ValidateFailed(ctx, "invalid-password", nil)
	case errors.Is(err, ErrInvalidCitizenID):
		base.ValidateFailed(ctx, "invalid-citizen-id", nil)
	case errors.Is(err, ErrInvalidPhone):
		base.ValidateFailed(ctx, "invalid-phone", nil)
	case errors.Is(err, ErrInvalidBirthDate):
		base.ValidateFailed(ctx, "invalid-birth-date", nil)
	case errors.Is(err, ErrInvalidNamePair):
		base.ValidateFailed(ctx, "invalid-name-pair", nil)
	case errors.Is(err, ErrInvalidApprovalReason):
		base.ValidateFailed(ctx, "invalid-approval-reason", nil)
	case errors.Is(err, ErrMemberStudentConditionFail):
		base.ValidateFailed(ctx, "condition-fail", nil)
	case errors.Is(err, ErrMemberStudentDuplicate):
		base.BadRequest(ctx, "member-student-duplicate", nil)
	default:
		log.Errf("%s: %v", fallback, err)
		base.InternalServerError(ctx, fallback, nil)
	}
}

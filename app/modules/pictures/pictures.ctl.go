package pictures

import (
	"errors"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/modules/s3"
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

func (c *Controller) isElevatedRole(role ent.MemberRole) bool {
	return role == ent.MemberRoleAdmin || role == ent.MemberRoleSuperadmin
}

func (c *Controller) handleServiceError(ctx *gin.Context, log *logpkg.Logger, err error, fallback string) {
	switch {
	case errors.Is(err, ErrPictureNotFound):
		base.BadRequest(ctx, "picture-not-found", nil)
	case errors.Is(err, ErrPictureDuplicate):
		base.BadRequest(ctx, "picture-duplicate", nil)
	case errors.Is(err, ErrPictureUnauthorized):
		base.Unauthorized(ctx, "unauthorized", nil)
	case errors.Is(err, ErrPictureConditionFail):
		base.ValidateFailed(ctx, "condition-fail", nil)
	case errors.Is(err, s3.ErrS3InvalidObjectKey):
		base.BadRequest(ctx, "invalid-picture-object-key", nil)
	default:
		log.Errf("%s: %v", fallback, err)
		base.InternalServerError(ctx, fallback, nil)
	}
}

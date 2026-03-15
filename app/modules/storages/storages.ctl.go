package storages

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
	case errors.Is(err, ErrStorageNotFound):
		base.BadRequest(ctx, "storage-not-found", nil)
	case errors.Is(err, ErrStorageConditionFail):
		base.ValidateFailed(ctx, "condition-fail", nil)
	case errors.Is(err, ErrStorageDuplicate):
		base.BadRequest(ctx, "storage-duplicate", nil)
	default:
		log.Errf("%s: %v", fallback, err)
		base.InternalServerError(ctx, fallback, nil)
	}
}

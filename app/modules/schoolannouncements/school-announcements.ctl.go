package schoolannouncements

import (
	"errors"

	"eduflow/app/modules/auth"
	"eduflow/app/modules/entities/ent"
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

func (c *Controller) currentUser(ctx *gin.Context) (*auth.CurrentUser, bool) {
	user, ok := auth.CurrentUserFromGin(ctx)
	if !ok || user == nil || user.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return nil, false
	}
	return user, true
}

func (c *Controller) canWrite(user *auth.CurrentUser) bool {
	if user == nil || user.Member == nil {
		return false
	}
	switch user.Member.Role {
	case ent.MemberRoleSuperadmin, ent.MemberRoleAdmin, ent.MemberRoleStaff:
		return true
	default:
		return false
	}
}

func (c *Controller) handleServiceError(ctx *gin.Context, log *logpkg.Logger, err error, fallback string) {
	switch {
	case errors.Is(err, ErrSchoolAnnouncementNotFound):
		base.BadRequest(ctx, "school-announcement-not-found", nil)
	case errors.Is(err, ErrSchoolAnnouncementUnauthorized):
		base.Unauthorized(ctx, "unauthorized", nil)
	case errors.Is(err, ErrSchoolAnnouncementConditionFail):
		base.ValidateFailed(ctx, "condition-fail", nil)
	case errors.Is(err, ErrSchoolAnnouncementInvalidUpdate):
		base.ValidateFailed(ctx, "condition-fail", nil)
	default:
		log.Errf("%s: %v", fallback, err)
		base.InternalServerError(ctx, fallback, nil)
	}
}

package reports

import (
	"eduflow/app/modules/auth"
	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"

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

func (c *Controller) canRead(user *auth.CurrentUser) bool {
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

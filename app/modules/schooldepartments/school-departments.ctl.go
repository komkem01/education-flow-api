package schooldepartments

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

func (c *Controller) handleServiceError(ctx *gin.Context, log *logpkg.Logger, err error, fallback string) {
	switch {
	case errors.Is(err, ErrSchoolDepartmentNotFound):
		base.BadRequest(ctx, "school-department-not-found", nil)
	case errors.Is(err, ErrSchoolDepartmentDuplicate):
		base.BadRequest(ctx, "school-department-duplicate", nil)
	case errors.Is(err, ErrSchoolDepartmentConditionFail):
		base.ValidateFailed(ctx, "condition-fail", nil)
	default:
		log.Errf("%s: %v", fallback, err)
		base.InternalServerError(ctx, fallback, nil)
	}
}

func canWriteSchoolDepartments(role ent.MemberRole) bool {
	switch role {
	case ent.MemberRoleSuperadmin, ent.MemberRoleAdmin, ent.MemberRoleStaff:
		return true
	default:
		return false
	}
}

package memberstudents

import (
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Module struct {
	tracer trace.Tracer
	Svc    *Service
	Ctl    *Controller
}

func New(conf *config.Config[Config], ent entitiesinf.MemberStudentEntity, approvalEnt entitiesinf.ApprovalRequestEntity, actionEnt entitiesinf.ApprovalActionEntity) *Module {
	tracer := otel.Tracer("eduflow.modules.memberstudents")
	svc := newService(&Options{Config: conf, tracer: tracer, db: ent, approval: approvalEnt, action: actionEnt})

	return &Module{tracer: tracer, Svc: svc, Ctl: newController(tracer, svc)}
}

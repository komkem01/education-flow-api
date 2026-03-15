package approvals

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

func New(
	conf *config.Config[Config],
	ent entitiesinf.ApprovalRequestEntity,
	actionEnt entitiesinf.ApprovalActionEntity,
	memberEnt entitiesinf.MemberEntity,
	studentEnt entitiesinf.MemberStudentEntity,
	managementEnt entitiesinf.MemberManagementEntity,
) *Module {
	tracer := otel.Tracer("eduflow.modules.approvals")
	svc := newService(&Options{
		Config:       conf,
		tracer:       tracer,
		db:           ent,
		actionDB:     actionEnt,
		memberDB:     memberEnt,
		studentDB:    studentEnt,
		managementDB: managementEnt,
	})

	return &Module{tracer: tracer, Svc: svc, Ctl: newController(tracer, svc)}
}

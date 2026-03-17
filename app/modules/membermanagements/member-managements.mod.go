package membermanagements

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

func New(conf *config.Config[Config], ent entitiesinf.MemberManagementEntity, schoolDepartmentEnt entitiesinf.SchoolDepartmentEntity) *Module {
	tracer := otel.Tracer("eduflow.modules.membermanagements")
	svc := newService(&Options{Config: conf, tracer: tracer, db: ent, schoolDepartmentDB: schoolDepartmentEnt})

	return &Module{tracer: tracer, Svc: svc, Ctl: newController(tracer, svc)}
}

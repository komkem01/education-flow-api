package teachereducations

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

func New(conf *config.Config[Config], educationEnt entitiesinf.TeacherEducationEntity) *Module {
	tracer := otel.Tracer("eduflow.modules.teachereducations")
	svc := newService(&Options{Config: conf, tracer: tracer, db: educationEnt})

	return &Module{tracer: tracer, Svc: svc, Ctl: newController(tracer, svc)}
}

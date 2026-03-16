package auth

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
	memberEnt entitiesinf.MemberEntity,
	sessionEnt entitiesinf.AuthSessionEntity,
	schoolEnt entitiesinf.SchoolEntity,
	memberTeacherEnt entitiesinf.MemberTeacherEntity,
) *Module {
	tracer := otel.Tracer("eduflow.modules.auth")
	svc := newService(&Options{
		Config:        conf,
		tracer:        tracer,
		member:        memberEnt,
		session:       sessionEnt,
		school:        schoolEnt,
		memberTeacher: memberTeacherEnt,
	})

	return &Module{tracer: tracer, Svc: svc, Ctl: newController(tracer, svc)}
}

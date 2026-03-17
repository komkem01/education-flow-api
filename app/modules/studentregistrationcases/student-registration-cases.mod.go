package studentregistrationcases

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

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.StudentRegistrationCaseEntity
}

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.StudentRegistrationCaseEntity
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func New(conf *config.Config[Config], registrationEnt entitiesinf.StudentRegistrationCaseEntity) *Module {
	tracer := otel.Tracer("eduflow.modules.studentregistrationcases")
	svc := newService(&Options{Config: conf, tracer: tracer, db: registrationEnt})
	return &Module{tracer: tracer, Svc: svc, Ctl: newController(tracer, svc)}
}

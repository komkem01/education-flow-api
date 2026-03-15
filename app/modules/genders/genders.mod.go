package genders

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

func New(conf *config.Config[Config], genderEnt entitiesinf.GenderEntity) *Module {
	tracer := otel.Tracer("eduflow.modules.genders")
	svc := newService(&Options{
		Config: conf,
		tracer: tracer,
		db:     genderEnt,
	})

	return &Module{
		tracer: tracer,
		Svc:    svc,
		Ctl:    newController(tracer, svc),
	}
}

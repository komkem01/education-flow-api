package s3

import (
	"eduflow/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Module struct {
	tracer trace.Tracer
	Svc    *Service
}

func New(conf *config.Config[Config]) *Module {
	tracer := otel.Tracer("eduflow.modules.s3")
	svc, err := newService(&Options{Config: conf, tracer: tracer})
	if err != nil {
		panic(err)
	}

	return &Module{tracer: tracer, Svc: svc}
}

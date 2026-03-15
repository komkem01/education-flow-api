package documents

import (
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/modules/s3"
	"eduflow/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Module struct {
	tracer trace.Tracer
	Svc    *Service
	Ctl    *Controller
}

func New(conf *config.Config[Config], s3Svc *s3.Service, db entitiesinf.DocumentEntity, storageDB entitiesinf.StorageEntity) *Module {
	tracer := otel.Tracer("eduflow.modules.documents")
	svc := newService(&Options{Config: conf, tracer: tracer, s3: s3Svc, db: db, storageDB: storageDB})

	return &Module{tracer: tracer, Svc: svc, Ctl: newController(tracer, svc)}
}

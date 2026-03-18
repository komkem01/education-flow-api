package storages

import (
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Module struct {
	Svc *Service
	Ctl *Controller
}

type Config struct{}

type RailwayConfig struct {
	URL            string
	ServiceRoleKey string
	PublicBucket   string
	PrivateBucket  string
}

type (
	Service struct {
		tracer         trace.Tracer
		db             entitiesinf.StorageEntity
		railwayStorage *railwayStorageClient
	}
)

type Options struct {
	*config.Config[Config]
	tracer      trace.Tracer
	db          entitiesinf.StorageEntity
	railwayConf RailwayConfig
}

func New(conf *config.Config[Config], db entitiesinf.StorageEntity) *Module {
	tracer := otel.Tracer("storages_module")
	svc := newService(&Options{
		Config:      conf,
		tracer:      tracer,
		db:          db,
		railwayConf: RailwayConfig{},
	})
	return &Module{Svc: svc, Ctl: newController(tracer, svc)}
}

func newService(opt *Options) *Service {
	return &Service{
		tracer:         opt.tracer,
		db:             opt.db,
		railwayStorage: newRailwayStorageClient(opt.railwayConf),
	}
}

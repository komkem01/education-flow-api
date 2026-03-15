package attendancesessions

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

func New(conf *config.Config[Config], ent entitiesinf.AttendanceSessionEntity, enrollment entitiesinf.StudentEnrollmentEntity, record entitiesinf.AttendanceRecordEntity) *Module {
	tracer := otel.Tracer("eduflow.modules.attendancesessions")
	svc := newService(&Options{Config: conf, tracer: tracer, db: ent, enrollment: enrollment, record: record})

	return &Module{tracer: tracer, Svc: svc, Ctl: newController(tracer, svc)}
}

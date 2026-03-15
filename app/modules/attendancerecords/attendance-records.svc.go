package attendancerecords

import (
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/internal/config"

	"go.opentelemetry.io/otel/trace"
)

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.AttendanceRecordEntity
	logDB  entitiesinf.AttendanceRecordLogEntity
}

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.AttendanceRecordEntity
	logDB  entitiesinf.AttendanceRecordLogEntity
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db, logDB: opt.logDB}
}

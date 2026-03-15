package approvals

import (
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/internal/config"

	"go.opentelemetry.io/otel/trace"
)

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer       trace.Tracer
	db           entitiesinf.ApprovalRequestEntity
	actionDB     entitiesinf.ApprovalActionEntity
	memberDB     entitiesinf.MemberEntity
	studentDB    entitiesinf.MemberStudentEntity
	managementDB entitiesinf.MemberManagementEntity
}

type Service struct {
	tracer       trace.Tracer
	db           entitiesinf.ApprovalRequestEntity
	actionDB     entitiesinf.ApprovalActionEntity
	memberDB     entitiesinf.MemberEntity
	studentDB    entitiesinf.MemberStudentEntity
	managementDB entitiesinf.MemberManagementEntity
}

func newService(opt *Options) *Service {
	return &Service{
		tracer:       opt.tracer,
		db:           opt.db,
		actionDB:     opt.actionDB,
		memberDB:     opt.memberDB,
		studentDB:    opt.studentDB,
		managementDB: opt.managementDB,
	}
}

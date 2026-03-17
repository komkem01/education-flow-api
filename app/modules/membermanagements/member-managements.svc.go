package membermanagements

import (
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/internal/config"

	"go.opentelemetry.io/otel/trace"
)

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer             trace.Tracer
	db                 entitiesinf.MemberManagementEntity
	schoolDepartmentDB entitiesinf.SchoolDepartmentEntity
	approvalDB         entitiesinf.ApprovalRequestEntity
	actionDB           entitiesinf.ApprovalActionEntity
}

type Service struct {
	tracer             trace.Tracer
	db                 entitiesinf.MemberManagementEntity
	schoolDepartmentDB entitiesinf.SchoolDepartmentEntity
	approvalDB         entitiesinf.ApprovalRequestEntity
	actionDB           entitiesinf.ApprovalActionEntity
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db, schoolDepartmentDB: opt.schoolDepartmentDB, approvalDB: opt.approvalDB, actionDB: opt.actionDB}
}

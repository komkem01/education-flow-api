package approvals

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
	ent entitiesinf.ApprovalRequestEntity,
	actionEnt entitiesinf.ApprovalActionEntity,
	studentRegEnt entitiesinf.StudentRegistrationCaseEntity,
	memberEnt entitiesinf.MemberEntity,
	studentEnt entitiesinf.MemberStudentEntity,
	profileEnt entitiesinf.StudentProfileEntity,
	healthEnt entitiesinf.StudentHealthProfileEntity,
	managementEnt entitiesinf.MemberManagementEntity,
	teacherEnt entitiesinf.MemberTeacherEntity,
) *Module {
	tracer := otel.Tracer("eduflow.modules.approvals")
	svc := newService(&Options{
		Config:       conf,
		tracer:       tracer,
		db:           ent,
		actionDB:     actionEnt,
		studentRegDB: studentRegEnt,
		memberDB:     memberEnt,
		studentDB:    studentEnt,
		profileDB:    profileEnt,
		healthDB:     healthEnt,
		managementDB: managementEnt,
		teacherDB:    teacherEnt,
	})

	return &Module{tracer: tracer, Svc: svc, Ctl: newController(tracer, svc)}
}

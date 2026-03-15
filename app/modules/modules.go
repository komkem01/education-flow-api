package modules

import (
	academicyears "eduflow/app/modules/academicyears"
	"eduflow/app/modules/classrooms"
	"eduflow/app/modules/departments"
	"log/slog"
	"sync"

	"eduflow/app/modules/entities"
	"eduflow/app/modules/example"
	"eduflow/app/modules/genders"
	"eduflow/app/modules/membermanagements"
	"eduflow/app/modules/members"
	"eduflow/app/modules/memberteachers"
	"eduflow/app/modules/prefixes"
	"eduflow/app/modules/schools"
	"eduflow/app/modules/sentry"
	"eduflow/app/modules/specs"
	"eduflow/app/modules/subjectgroups"
	"eduflow/app/modules/subjects"
	"eduflow/app/modules/teachereducations"
	"eduflow/app/modules/teacherexperiences"
	"eduflow/app/modules/teacherrequests"
	"eduflow/internal/config"
	"eduflow/internal/database"
	"eduflow/internal/log"
	"eduflow/internal/otel/collector"

	exampletwo "eduflow/app/modules/example-two"

	appConf "eduflow/config"
	// "eduflow/app/modules/kafka"
)

type Modules struct {
	Conf   *config.Module[appConf.Config]
	Specs  *specs.Module
	Log    *log.Module
	OTEL   *collector.Module
	Sentry *sentry.Module
	DB     *database.DatabaseModule
	ENT    *entities.Module
	// Kafka *kafka.Module
	Example            *example.Module
	Example2           *exampletwo.Module
	Genders            *genders.Module
	Prefixes           *prefixes.Module
	Schools            *schools.Module
	Departments        *departments.Module
	AcademicYears      *academicyears.Module
	Classrooms         *classrooms.Module
	SubjectGroups      *subjectgroups.Module
	Subjects           *subjects.Module
	Members            *members.Module
	MemberManagements  *membermanagements.Module
	MemberTeachers     *memberteachers.Module
	TeacherEducations  *teachereducations.Module
	TeacherExperiences *teacherexperiences.Module
	TeacherRequests    *teacherrequests.Module
}

func modulesInit() {
	confMod := config.New(&appConf.App)
	specsMod := specs.New(config.Conf[specs.Config](confMod.Svc))
	conf := confMod.Svc.Config()

	logMod := log.New(config.Conf[log.Option](confMod.Svc))
	otel := collector.New(config.Conf[collector.Config](confMod.Svc))
	log := log.With(slog.String("module", "modules"))

	sentryMod := sentry.New(config.Conf[sentry.Config](confMod.Svc))

	db := database.New(conf.Database.Sql)
	entitiesMod := entities.New(db.Svc.DB())
	exampleMod := example.New(config.Conf[example.Config](confMod.Svc), entitiesMod.Svc)
	exampleMod2 := exampletwo.New(config.Conf[exampletwo.Config](confMod.Svc), entitiesMod.Svc)
	gendersMod := genders.New(config.Conf[genders.Config](confMod.Svc), entitiesMod.Svc)
	prefixesMod := prefixes.New(config.Conf[prefixes.Config](confMod.Svc), entitiesMod.Svc)
	schoolsMod := schools.New(config.Conf[schools.Config](confMod.Svc), entitiesMod.Svc)
	departmentsMod := departments.New(config.Conf[departments.Config](confMod.Svc), entitiesMod.Svc)
	academicYearsMod := academicyears.New(config.Conf[academicyears.Config](confMod.Svc), entitiesMod.Svc)
	classroomsMod := classrooms.New(config.Conf[classrooms.Config](confMod.Svc), entitiesMod.Svc)
	subjectGroupsMod := subjectgroups.New(config.Conf[subjectgroups.Config](confMod.Svc), entitiesMod.Svc)
	subjectsMod := subjects.New(config.Conf[subjects.Config](confMod.Svc), entitiesMod.Svc)
	membersMod := members.New(config.Conf[members.Config](confMod.Svc), entitiesMod.Svc)
	memberManagementsMod := membermanagements.New(config.Conf[membermanagements.Config](confMod.Svc), entitiesMod.Svc)
	memberTeachersMod := memberteachers.New(config.Conf[memberteachers.Config](confMod.Svc), entitiesMod.Svc)
	teacherEducationsMod := teachereducations.New(config.Conf[teachereducations.Config](confMod.Svc), entitiesMod.Svc)
	teacherExperiencesMod := teacherexperiences.New(config.Conf[teacherexperiences.Config](confMod.Svc), entitiesMod.Svc)
	teacherRequestsMod := teacherrequests.New(config.Conf[teacherrequests.Config](confMod.Svc), entitiesMod.Svc)
	// kafka := kafka.New(&conf.Kafka)
	mod = &Modules{
		Conf:               confMod,
		Specs:              specsMod,
		Log:                logMod,
		OTEL:               otel,
		Sentry:             sentryMod,
		DB:                 db,
		ENT:                entitiesMod,
		Example:            exampleMod,
		Example2:           exampleMod2,
		Genders:            gendersMod,
		Prefixes:           prefixesMod,
		Schools:            schoolsMod,
		Departments:        departmentsMod,
		AcademicYears:      academicYearsMod,
		Classrooms:         classroomsMod,
		SubjectGroups:      subjectGroupsMod,
		Subjects:           subjectsMod,
		Members:            membersMod,
		MemberManagements:  memberManagementsMod,
		MemberTeachers:     memberTeachersMod,
		TeacherEducations:  teacherEducationsMod,
		TeacherExperiences: teacherExperiencesMod,
		TeacherRequests:    teacherRequestsMod,
	}

	log.Infof("all modules initialized")
}

var (
	once sync.Once
	mod  *Modules
)

func Get() *Modules {
	once.Do(modulesInit)

	return mod
}

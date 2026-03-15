package config

import (
	academicyears "eduflow/app/modules/academicyears"
	"eduflow/app/modules/departments"
	"eduflow/app/modules/example"
	exampletwo "eduflow/app/modules/example-two"
	"eduflow/app/modules/genders"
	"eduflow/app/modules/members"
	"eduflow/app/modules/memberteachers"
	"eduflow/app/modules/prefixes"
	"eduflow/app/modules/schools"
	"eduflow/app/modules/sentry"
	"eduflow/app/modules/specs"
	"eduflow/app/modules/teachereducations"
	"eduflow/app/modules/teacherexperiences"
	"eduflow/app/modules/teacherrequests"
	"eduflow/internal/kafka"
	"eduflow/internal/log"
	"eduflow/internal/otel/collector"
)

// Config is a struct that contains all the configuration of the application.
type Config struct {
	Database Database

	AppName     string
	AppKey      string
	Environment string
	Specs       specs.Config
	Debug       bool

	Port           int
	HttpJsonNaming string

	SslCaPath      string
	SslPrivatePath string
	SslCertPath    string

	Otel   collector.Config
	Sentry sentry.Config

	Kafka kafka.Config
	Log   log.Option

	Example example.Config

	ExampleTwo         exampletwo.Config
	Genders            genders.Config
	Prefixes           prefixes.Config
	Schools            schools.Config
	Departments        departments.Config
	AcademicYears      academicyears.Config
	Members            members.Config
	MemberTeachers     memberteachers.Config
	TeacherEducations  teachereducations.Config
	TeacherExperiences teacherexperiences.Config
	TeacherRequests    teacherrequests.Config
}

var App = Config{
	Specs: specs.Config{
		Version: "v1",
	},
	Database: database,
	Kafka:    kafkaConf,

	AppName: "go_app",
	Port:    8080,
	AppKey:  "secret",
	Debug:   false,

	HttpJsonNaming: "snake_case",

	SslCaPath:      "eduflow/cert/ca.pem",
	SslPrivatePath: "eduflow/cert/server.pem",
	SslCertPath:    "eduflow/cert/server-key.pem",

	Otel: collector.Config{
		CollectorEndpoint: "",
		LogMode:           "noop",
		TraceMode:         "noop",
		MetricMode:        "noop",
		TraceRatio:        0.01,
	},
}

package config

import (
	academicyears "eduflow/app/modules/academicyears"
	"eduflow/app/modules/approvals"
	"eduflow/app/modules/attendancerecordlogs"
	"eduflow/app/modules/attendancerecords"
	"eduflow/app/modules/attendancesessions"
	"eduflow/app/modules/auditlogs"
	"eduflow/app/modules/auth"
	"eduflow/app/modules/classrooms"
	"eduflow/app/modules/departments"
	"eduflow/app/modules/documents"
	"eduflow/app/modules/enrollmentstatushistories"
	"eduflow/app/modules/enrollmentsubjects"
	"eduflow/app/modules/example"
	exampletwo "eduflow/app/modules/example-two"
	"eduflow/app/modules/genders"
	"eduflow/app/modules/memberguardians"
	"eduflow/app/modules/membermanagements"
	"eduflow/app/modules/members"
	"eduflow/app/modules/memberstudents"
	"eduflow/app/modules/memberteachers"
	"eduflow/app/modules/pictures"
	"eduflow/app/modules/prefixes"
	"eduflow/app/modules/reports"
	"eduflow/app/modules/s3"
	"eduflow/app/modules/schoolannouncements"
	"eduflow/app/modules/schooldepartments"
	"eduflow/app/modules/schools"
	"eduflow/app/modules/sentry"
	"eduflow/app/modules/specs"
	"eduflow/app/modules/storages"
	"eduflow/app/modules/studentenrollments"
	"eduflow/app/modules/studentguardians"
	"eduflow/app/modules/studenthealthprofiles"
	"eduflow/app/modules/studentprofiles"
	"eduflow/app/modules/studentregistrationcases"
	"eduflow/app/modules/subjectgroups"
	"eduflow/app/modules/subjects"
	"eduflow/app/modules/teachereducations"
	"eduflow/app/modules/teacheremergencycontacts"
	"eduflow/app/modules/teacherexperiences"
	"eduflow/app/modules/teacherhealthprofiles"
	"eduflow/app/modules/teacherlicenses"
	"eduflow/app/modules/teacherrequests"
	"eduflow/app/modules/teachersubjects"
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

	ExampleTwo                exampletwo.Config
	Genders                   genders.Config
	Prefixes                  prefixes.Config
	Schools                   schools.Config
	Departments               departments.Config
	AcademicYears             academicyears.Config
	Classrooms                classrooms.Config
	SubjectGroups             subjectgroups.Config
	Subjects                  subjects.Config
	SchoolAnnouncements       schoolannouncements.Config
	AuditLogs                 auditlogs.Config
	S3                        s3.Config
	Storages                  storages.Config
	Documents                 documents.Config
	Pictures                  pictures.Config
	Auth                      auth.Config
	Approvals                 approvals.Config
	AttendanceSessions        attendancesessions.Config
	AttendanceRecords         attendancerecords.Config
	AttendanceRecordLogs      attendancerecordlogs.Config
	Members                   members.Config
	MemberStudents            memberstudents.Config
	MemberGuardians           memberguardians.Config
	StudentGuardians          studentguardians.Config
	StudentEnrollments        studentenrollments.Config
	EnrollmentStatusHistories enrollmentstatushistories.Config
	EnrollmentSubjects        enrollmentsubjects.Config
	StudentProfiles           studentprofiles.Config
	StudentHealthProfiles     studenthealthprofiles.Config
	StudentRegistrationCases  studentregistrationcases.Config
	MemberManagements         membermanagements.Config
	MemberTeachers            memberteachers.Config
	Reports                   reports.Config
	SchoolDepartments         schooldepartments.Config
	TeacherEducations         teachereducations.Config
	TeacherEmergencyContacts  teacheremergencycontacts.Config
	TeacherLicenses           teacherlicenses.Config
	TeacherExperiences        teacherexperiences.Config
	TeacherHealthProfiles     teacherhealthprofiles.Config
	TeacherRequests           teacherrequests.Config
	TeacherSubjects           teachersubjects.Config
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

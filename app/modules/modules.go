package modules

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
	"log/slog"
	"sync"

	"eduflow/app/modules/enrollmentstatushistories"
	"eduflow/app/modules/enrollmentsubjects"
	"eduflow/app/modules/entities"
	"eduflow/app/modules/example"
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
	Example                   *example.Module
	Example2                  *exampletwo.Module
	Genders                   *genders.Module
	Prefixes                  *prefixes.Module
	Schools                   *schools.Module
	Departments               *departments.Module
	AcademicYears             *academicyears.Module
	Classrooms                *classrooms.Module
	SubjectGroups             *subjectgroups.Module
	Subjects                  *subjects.Module
	SchoolAnnouncements       *schoolannouncements.Module
	AuditLogs                 *auditlogs.Module
	S3                        *s3.Module
	Storages                  *storages.Module
	Documents                 *documents.Module
	Pictures                  *pictures.Module
	Auth                      *auth.Module
	Approvals                 *approvals.Module
	AttendanceSessions        *attendancesessions.Module
	AttendanceRecords         *attendancerecords.Module
	AttendanceRecordLogs      *attendancerecordlogs.Module
	Members                   *members.Module
	MemberStudents            *memberstudents.Module
	MemberGuardians           *memberguardians.Module
	StudentGuardians          *studentguardians.Module
	StudentEnrollments        *studentenrollments.Module
	EnrollmentStatusHistories *enrollmentstatushistories.Module
	EnrollmentSubjects        *enrollmentsubjects.Module
	StudentProfiles           *studentprofiles.Module
	StudentHealthProfiles     *studenthealthprofiles.Module
	StudentRegistrationCases  *studentregistrationcases.Module
	MemberManagements         *membermanagements.Module
	MemberTeachers            *memberteachers.Module
	Reports                   *reports.Module
	SchoolDepartments         *schooldepartments.Module
	TeacherEducations         *teachereducations.Module
	TeacherEmergencyContacts  *teacheremergencycontacts.Module
	TeacherLicenses           *teacherlicenses.Module
	TeacherExperiences        *teacherexperiences.Module
	TeacherHealthProfiles     *teacherhealthprofiles.Module
	TeacherRequests           *teacherrequests.Module
	TeacherSubjects           *teachersubjects.Module
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
	schoolsMod := schools.New(config.Conf[schools.Config](confMod.Svc), entitiesMod.Svc, entitiesMod.Svc)
	departmentsMod := departments.New(config.Conf[departments.Config](confMod.Svc), entitiesMod.Svc)
	academicYearsMod := academicyears.New(config.Conf[academicyears.Config](confMod.Svc), entitiesMod.Svc)
	classroomsMod := classrooms.New(config.Conf[classrooms.Config](confMod.Svc), entitiesMod.Svc)
	subjectGroupsMod := subjectgroups.New(config.Conf[subjectgroups.Config](confMod.Svc), entitiesMod.Svc)
	subjectsMod := subjects.New(config.Conf[subjects.Config](confMod.Svc), entitiesMod.Svc)
	schoolAnnouncementsMod := schoolannouncements.New(config.Conf[schoolannouncements.Config](confMod.Svc), entitiesMod.Svc)
	auditLogsMod := auditlogs.New(config.Conf[auditlogs.Config](confMod.Svc), entitiesMod.Svc)
	s3Mod := s3.New(config.Conf[s3.Config](confMod.Svc))
	storagesMod := storages.New(config.Conf[storages.Config](confMod.Svc), entitiesMod.Svc)
	documentsMod := documents.New(config.Conf[documents.Config](confMod.Svc), s3Mod.Svc, entitiesMod.Svc, entitiesMod.Svc)
	picturesMod := pictures.New(config.Conf[pictures.Config](confMod.Svc), s3Mod.Svc, entitiesMod.Svc, entitiesMod.Svc)
	authMod := auth.New(config.Conf[auth.Config](confMod.Svc), entitiesMod.Svc, entitiesMod.Svc, entitiesMod.Svc, entitiesMod.Svc)
	approvalsMod := approvals.New(
		config.Conf[approvals.Config](confMod.Svc),
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
		entitiesMod.Svc,
	)
	attendanceSessionsMod := attendancesessions.New(config.Conf[attendancesessions.Config](confMod.Svc), entitiesMod.Svc, entitiesMod.Svc, entitiesMod.Svc)
	attendanceRecordsMod := attendancerecords.New(config.Conf[attendancerecords.Config](confMod.Svc), entitiesMod.Svc, entitiesMod.Svc)
	attendanceRecordLogsMod := attendancerecordlogs.New(config.Conf[attendancerecordlogs.Config](confMod.Svc), entitiesMod.Svc)
	membersMod := members.New(config.Conf[members.Config](confMod.Svc), entitiesMod.Svc)
	memberStudentsMod := memberstudents.New(config.Conf[memberstudents.Config](confMod.Svc), entitiesMod.Svc, entitiesMod.Svc, entitiesMod.Svc)
	memberGuardiansMod := memberguardians.New(config.Conf[memberguardians.Config](confMod.Svc), entitiesMod.Svc)
	studentGuardiansMod := studentguardians.New(config.Conf[studentguardians.Config](confMod.Svc), entitiesMod.Svc)
	studentEnrollmentsMod := studentenrollments.New(config.Conf[studentenrollments.Config](confMod.Svc), entitiesMod.Svc)
	enrollmentStatusHistoriesMod := enrollmentstatushistories.New(config.Conf[enrollmentstatushistories.Config](confMod.Svc), entitiesMod.Svc)
	enrollmentSubjectsMod := enrollmentsubjects.New(config.Conf[enrollmentsubjects.Config](confMod.Svc), entitiesMod.Svc)
	studentProfilesMod := studentprofiles.New(config.Conf[studentprofiles.Config](confMod.Svc), entitiesMod.Svc)
	studentHealthProfilesMod := studenthealthprofiles.New(config.Conf[studenthealthprofiles.Config](confMod.Svc), entitiesMod.Svc)
	studentRegistrationCasesMod := studentregistrationcases.New(config.Conf[studentregistrationcases.Config](confMod.Svc), entitiesMod.Svc)
	memberManagementsMod := membermanagements.New(config.Conf[membermanagements.Config](confMod.Svc), entitiesMod.Svc, entitiesMod.Svc, entitiesMod.Svc, entitiesMod.Svc)
	memberTeachersMod := memberteachers.New(config.Conf[memberteachers.Config](confMod.Svc), entitiesMod.Svc, entitiesMod.Svc, entitiesMod.Svc)
	reportsMod := reports.New(config.Conf[reports.Config](confMod.Svc), entitiesMod.Svc)
	schoolDepartmentsMod := schooldepartments.New(config.Conf[schooldepartments.Config](confMod.Svc), entitiesMod.Svc)
	teacherEducationsMod := teachereducations.New(config.Conf[teachereducations.Config](confMod.Svc), entitiesMod.Svc)
	teacherEmergencyContactsMod := teacheremergencycontacts.New(config.Conf[teacheremergencycontacts.Config](confMod.Svc), entitiesMod.Svc)
	teacherLicensesMod := teacherlicenses.New(config.Conf[teacherlicenses.Config](confMod.Svc), entitiesMod.Svc)
	teacherExperiencesMod := teacherexperiences.New(config.Conf[teacherexperiences.Config](confMod.Svc), entitiesMod.Svc)
	teacherHealthProfilesMod := teacherhealthprofiles.New(config.Conf[teacherhealthprofiles.Config](confMod.Svc), entitiesMod.Svc)
	teacherRequestsMod := teacherrequests.New(config.Conf[teacherrequests.Config](confMod.Svc), entitiesMod.Svc, entitiesMod.Svc, entitiesMod.Svc)
	teacherSubjectsMod := teachersubjects.New(config.Conf[teachersubjects.Config](confMod.Svc), entitiesMod.Svc)
	// kafka := kafka.New(&conf.Kafka)
	mod = &Modules{
		Conf:                      confMod,
		Specs:                     specsMod,
		Log:                       logMod,
		OTEL:                      otel,
		Sentry:                    sentryMod,
		DB:                        db,
		ENT:                       entitiesMod,
		Example:                   exampleMod,
		Example2:                  exampleMod2,
		Genders:                   gendersMod,
		Prefixes:                  prefixesMod,
		Schools:                   schoolsMod,
		Departments:               departmentsMod,
		AcademicYears:             academicYearsMod,
		Classrooms:                classroomsMod,
		SubjectGroups:             subjectGroupsMod,
		Subjects:                  subjectsMod,
		SchoolAnnouncements:       schoolAnnouncementsMod,
		AuditLogs:                 auditLogsMod,
		S3:                        s3Mod,
		Storages:                  storagesMod,
		Documents:                 documentsMod,
		Pictures:                  picturesMod,
		Auth:                      authMod,
		Approvals:                 approvalsMod,
		AttendanceSessions:        attendanceSessionsMod,
		AttendanceRecords:         attendanceRecordsMod,
		AttendanceRecordLogs:      attendanceRecordLogsMod,
		Members:                   membersMod,
		MemberStudents:            memberStudentsMod,
		MemberGuardians:           memberGuardiansMod,
		StudentGuardians:          studentGuardiansMod,
		StudentEnrollments:        studentEnrollmentsMod,
		EnrollmentStatusHistories: enrollmentStatusHistoriesMod,
		EnrollmentSubjects:        enrollmentSubjectsMod,
		StudentProfiles:           studentProfilesMod,
		StudentHealthProfiles:     studentHealthProfilesMod,
		StudentRegistrationCases:  studentRegistrationCasesMod,
		MemberManagements:         memberManagementsMod,
		MemberTeachers:            memberTeachersMod,
		Reports:                   reportsMod,
		SchoolDepartments:         schoolDepartmentsMod,
		TeacherEducations:         teacherEducationsMod,
		TeacherEmergencyContacts:  teacherEmergencyContactsMod,
		TeacherLicenses:           teacherLicensesMod,
		TeacherExperiences:        teacherExperiencesMod,
		TeacherHealthProfiles:     teacherHealthProfilesMod,
		TeacherRequests:           teacherRequestsMod,
		TeacherSubjects:           teacherSubjectsMod,
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

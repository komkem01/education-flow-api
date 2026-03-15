package routes

import (
	"fmt"
	"net/http"

	"eduflow/app/modules"

	"github.com/gin-gonic/gin"
)

func WarpH(router *gin.RouterGroup, prefix string, handler http.Handler) {
	router.Any(fmt.Sprintf("%s/*w", prefix), gin.WrapH(http.StripPrefix(fmt.Sprintf("%s%s", router.BasePath(), prefix), handler)))
}

func api(r *gin.RouterGroup, mod *modules.Modules) {
	r.GET("/example/:id", mod.Example.Ctl.Get)
	r.GET("/example-http", mod.Example.Ctl.GetHttpReq)
	r.POST("/example", mod.Example.Ctl.Create)
}

func apiSystem(r *gin.RouterGroup, mod *modules.Modules) {
	// Public routes (no authentication required)
	system := r.Group("/system")
	{
		genders := system.Group("/genders")
		{
			genders.GET("", mod.Genders.Ctl.GendersList)
			genders.GET("/:id", mod.Genders.Ctl.GendersInfo)
			genders.POST("", mod.Genders.Ctl.CreateGenderController)
			genders.PATCH("/:id", mod.Genders.Ctl.GendersUpdate)
			genders.DELETE("/:id", mod.Genders.Ctl.GendersDelete)
		}
		prefixes := system.Group("/prefixes")
		{
			prefixes.GET("", mod.Prefixes.Ctl.PrefixesList)
			prefixes.GET("/:id", mod.Prefixes.Ctl.PrefixesInfo)
			prefixes.POST("", mod.Prefixes.Ctl.CreatePrefixController)
			prefixes.PATCH("/:id", mod.Prefixes.Ctl.PrefixesUpdate)
			prefixes.DELETE("/:id", mod.Prefixes.Ctl.PrefixesDelete)
		}
		schools := system.Group("/schools")
		{
			schools.GET("", mod.Schools.Ctl.SchoolsList)
			schools.GET("/:id", mod.Schools.Ctl.SchoolsInfo)
			schools.POST("", mod.Schools.Ctl.CreateSchoolController)
			schools.PATCH("/:id", mod.Schools.Ctl.SchoolsUpdate)
			schools.DELETE("/:id", mod.Schools.Ctl.SchoolsDelete)
		}
		academicYears := system.Group("/academic-years")
		{
			academicYears.GET("", mod.AcademicYears.Ctl.AcademicYearsList)
			academicYears.GET("/:id", mod.AcademicYears.Ctl.AcademicYearsInfo)
			academicYears.POST("", mod.AcademicYears.Ctl.CreateAcademicYearController)
			academicYears.PATCH("/:id", mod.AcademicYears.Ctl.AcademicYearsUpdate)
			academicYears.DELETE("/:id", mod.AcademicYears.Ctl.AcademicYearsDelete)
		}
		departments := system.Group("/departments")
		{
			departments.GET("", mod.Departments.Ctl.DepartmentsList)
			departments.GET("/:id", mod.Departments.Ctl.DepartmentsInfo)
			departments.POST("", mod.Departments.Ctl.CreateDepartmentController)
			departments.PATCH("/:id", mod.Departments.Ctl.DepartmentsUpdate)
			departments.DELETE("/:id", mod.Departments.Ctl.DepartmentsDelete)
		}
	}
}

func apiPublic(r *gin.RouterGroup, mod *modules.Modules) {
	// Public routes (no authentication required) e.g. for login, registration, etc.
	public := r.Group("/public")
	{
		public.POST("/example", mod.Example.Ctl.Create)
	}
}

func apiMember(r *gin.RouterGroup, mod *modules.Modules) {
	// Protected routes (authentication required)
	members := r.Group("/member")
	{

		members.GET("", mod.Members.Ctl.MembersList)
		members.GET("/:id", mod.Members.Ctl.MembersInfo)
		members.POST("", mod.Members.Ctl.CreateMemberController)
		members.PATCH("/:id", mod.Members.Ctl.MembersUpdate)
		members.DELETE("/:id", mod.Members.Ctl.MembersDelete)

		students := members.Group("/students")
		{
			students.GET("", mod.MemberStudents.Ctl.MemberStudentsList)
			students.GET("/:id", mod.MemberStudents.Ctl.MemberStudentsInfo)
			students.POST("", mod.MemberStudents.Ctl.CreateMemberStudentController)
			students.PATCH("/:id", mod.MemberStudents.Ctl.MemberStudentsUpdate)
			students.DELETE("/:id", mod.MemberStudents.Ctl.MemberStudentsDelete)
		}

		guardians := members.Group("/guardians")
		{
			guardians.GET("", mod.MemberGuardians.Ctl.MemberGuardiansList)
			guardians.GET("/:id", mod.MemberGuardians.Ctl.MemberGuardiansInfo)
			guardians.POST("", mod.MemberGuardians.Ctl.CreateMemberGuardianController)
			guardians.PATCH("/:id", mod.MemberGuardians.Ctl.MemberGuardiansUpdate)
			guardians.DELETE("/:id", mod.MemberGuardians.Ctl.MemberGuardiansDelete)
		}

		studentGuardians := members.Group("/student-guardians")
		{
			studentGuardians.GET("", mod.StudentGuardians.Ctl.StudentGuardiansList)
			studentGuardians.GET("/:id", mod.StudentGuardians.Ctl.StudentGuardiansInfo)
			studentGuardians.POST("", mod.StudentGuardians.Ctl.CreateStudentGuardianController)
			studentGuardians.PATCH("/:id", mod.StudentGuardians.Ctl.StudentGuardiansUpdate)
			studentGuardians.DELETE("/:id", mod.StudentGuardians.Ctl.StudentGuardiansDelete)
		}

		studentProfiles := members.Group("/student-profiles")
		{
			studentProfiles.GET("", mod.StudentProfiles.Ctl.StudentProfilesList)
			studentProfiles.GET("/:id", mod.StudentProfiles.Ctl.StudentProfilesInfo)
			studentProfiles.POST("", mod.StudentProfiles.Ctl.CreateStudentProfileController)
			studentProfiles.PATCH("/:id", mod.StudentProfiles.Ctl.StudentProfilesUpdate)
			studentProfiles.DELETE("/:id", mod.StudentProfiles.Ctl.StudentProfilesDelete)
		}

		studentHealthProfiles := members.Group("/student-health-profiles")
		{
			studentHealthProfiles.GET("", mod.StudentHealthProfiles.Ctl.StudentHealthProfilesList)
			studentHealthProfiles.GET("/:id", mod.StudentHealthProfiles.Ctl.StudentHealthProfilesInfo)
			studentHealthProfiles.POST("", mod.StudentHealthProfiles.Ctl.CreateStudentHealthProfileController)
			studentHealthProfiles.PATCH("/:id", mod.StudentHealthProfiles.Ctl.StudentHealthProfilesUpdate)
			studentHealthProfiles.DELETE("/:id", mod.StudentHealthProfiles.Ctl.StudentHealthProfilesDelete)
		}

		managements := members.Group("/managements")
		{
			managements.GET("", mod.MemberManagements.Ctl.MemberManagementsList)
			managements.GET("/:id", mod.MemberManagements.Ctl.MemberManagementsInfo)
			managements.POST("", mod.MemberManagements.Ctl.CreateMemberManagementController)
			managements.PATCH("/:id", mod.MemberManagements.Ctl.MemberManagementsUpdate)
			managements.DELETE("/:id", mod.MemberManagements.Ctl.MemberManagementsDelete)
		}

		teachers := members.Group("/teachers")
		{
			teachers.GET("", mod.MemberTeachers.Ctl.MemberTeachersList)
			teachers.GET("/:id", mod.MemberTeachers.Ctl.MemberTeachersInfo)
			teachers.POST("", mod.MemberTeachers.Ctl.CreateMemberTeacherController)
			teachers.PATCH("/:id", mod.MemberTeachers.Ctl.MemberTeachersUpdate)
			teachers.DELETE("/:id", mod.MemberTeachers.Ctl.MemberTeachersDelete)
		}

		educations := members.Group("/teacher-educations")
		{
			educations.GET("", mod.TeacherEducations.Ctl.TeacherEducationsList)
			educations.GET("/:id", mod.TeacherEducations.Ctl.TeacherEducationsInfo)
			educations.POST("", mod.TeacherEducations.Ctl.CreateTeacherEducationController)
			educations.PATCH("/:id", mod.TeacherEducations.Ctl.TeacherEducationsUpdate)
			educations.DELETE("/:id", mod.TeacherEducations.Ctl.TeacherEducationsDelete)
		}

		experiences := members.Group("/teacher-experiences")
		{
			experiences.GET("", mod.TeacherExperiences.Ctl.TeacherExperiencesList)
			experiences.GET("/:id", mod.TeacherExperiences.Ctl.TeacherExperiencesInfo)
			experiences.POST("", mod.TeacherExperiences.Ctl.CreateTeacherExperienceController)
			experiences.PATCH("/:id", mod.TeacherExperiences.Ctl.TeacherExperiencesUpdate)
			experiences.DELETE("/:id", mod.TeacherExperiences.Ctl.TeacherExperiencesDelete)
		}

		requests := members.Group("/teacher-requests")
		{
			requests.GET("", mod.TeacherRequests.Ctl.TeacherRequestsList)
			requests.GET("/:id", mod.TeacherRequests.Ctl.TeacherRequestsInfo)
			requests.POST("", mod.TeacherRequests.Ctl.CreateTeacherRequestController)
			requests.PATCH("/:id", mod.TeacherRequests.Ctl.TeacherRequestsUpdate)
			requests.DELETE("/:id", mod.TeacherRequests.Ctl.TeacherRequestsDelete)
		}

		licenses := members.Group("/teacher-licenses")
		{
			licenses.GET("", mod.TeacherLicenses.Ctl.TeacherLicensesList)
			licenses.GET("/:id", mod.TeacherLicenses.Ctl.TeacherLicensesInfo)
			licenses.POST("", mod.TeacherLicenses.Ctl.CreateTeacherLicenseController)
			licenses.PATCH("/:id", mod.TeacherLicenses.Ctl.TeacherLicensesUpdate)
			licenses.DELETE("/:id", mod.TeacherLicenses.Ctl.TeacherLicensesDelete)
		}

		subjects := members.Group("/teacher-subjects")
		{
			subjects.GET("", mod.TeacherSubjects.Ctl.TeacherSubjectsList)
			subjects.GET("/:id", mod.TeacherSubjects.Ctl.TeacherSubjectsInfo)
			subjects.POST("", mod.TeacherSubjects.Ctl.CreateTeacherSubjectController)
			subjects.PATCH("/:id", mod.TeacherSubjects.Ctl.TeacherSubjectsUpdate)
			subjects.DELETE("/:id", mod.TeacherSubjects.Ctl.TeacherSubjectsDelete)
		}

		healthProfiles := members.Group("/teacher-health-profiles")
		{
			healthProfiles.GET("", mod.TeacherHealthProfiles.Ctl.TeacherHealthProfilesList)
			healthProfiles.GET("/:id", mod.TeacherHealthProfiles.Ctl.TeacherHealthProfilesInfo)
			healthProfiles.POST("", mod.TeacherHealthProfiles.Ctl.CreateTeacherHealthProfileController)
			healthProfiles.PATCH("/:id", mod.TeacherHealthProfiles.Ctl.TeacherHealthProfilesUpdate)
			healthProfiles.DELETE("/:id", mod.TeacherHealthProfiles.Ctl.TeacherHealthProfilesDelete)
		}

		emergencyContacts := members.Group("/teacher-emergency-contacts")
		{
			emergencyContacts.GET("", mod.TeacherEmergencyContacts.Ctl.TeacherEmergencyContactsList)
			emergencyContacts.GET("/:id", mod.TeacherEmergencyContacts.Ctl.TeacherEmergencyContactsInfo)
			emergencyContacts.POST("", mod.TeacherEmergencyContacts.Ctl.CreateTeacherEmergencyContactController)
			emergencyContacts.PATCH("/:id", mod.TeacherEmergencyContacts.Ctl.TeacherEmergencyContactsUpdate)
			emergencyContacts.DELETE("/:id", mod.TeacherEmergencyContacts.Ctl.TeacherEmergencyContactsDelete)
		}
	}
}

func apiManagement(r *gin.RouterGroup, mod *modules.Modules) {
	// Management routes (authentication required) e.g. for administrative tasks, etc.
	management := r.Group("/management")
	{
		// Add management routes here
		classrooms := management.Group("/classrooms")
		{
			classrooms.GET("", mod.Classrooms.Ctl.ClassroomsList)
			classrooms.GET("/:id", mod.Classrooms.Ctl.ClassroomsInfo)
			classrooms.POST("", mod.Classrooms.Ctl.CreateClassroomController)
			classrooms.PATCH("/:id", mod.Classrooms.Ctl.ClassroomsUpdate)
			classrooms.DELETE("/:id", mod.Classrooms.Ctl.ClassroomsDelete)
		}

		subjectGroups := management.Group("/subject-groups")
		{
			subjectGroups.GET("", mod.SubjectGroups.Ctl.SubjectGroupsList)
			subjectGroups.GET("/:id", mod.SubjectGroups.Ctl.SubjectGroupsInfo)
			subjectGroups.POST("", mod.SubjectGroups.Ctl.CreateSubjectGroupController)
			subjectGroups.PATCH("/:id", mod.SubjectGroups.Ctl.SubjectGroupsUpdate)
			subjectGroups.DELETE("/:id", mod.SubjectGroups.Ctl.SubjectGroupsDelete)
		}

		subjects := management.Group("/subjects")
		{
			subjects.GET("", mod.Subjects.Ctl.SubjectsList)
			subjects.GET("/:id", mod.Subjects.Ctl.SubjectsInfo)
			subjects.POST("", mod.Subjects.Ctl.CreateSubjectController)
			subjects.PATCH("/:id", mod.Subjects.Ctl.SubjectsUpdate)
			subjects.DELETE("/:id", mod.Subjects.Ctl.SubjectsDelete)
		}
	}
}

func apiEnrollments(r *gin.RouterGroup, mod *modules.Modules) {
	// Enrollment routes (authentication required) e.g. for student enrollments, etc.
	enrollments := r.Group("/enrollments")
	{
		studentEnrollments := enrollments.Group("/student-enrollments")
		{
			studentEnrollments.GET("", mod.StudentEnrollments.Ctl.StudentEnrollmentsList)
			studentEnrollments.GET("/:id", mod.StudentEnrollments.Ctl.StudentEnrollmentsInfo)
			studentEnrollments.POST("", mod.StudentEnrollments.Ctl.CreateStudentEnrollmentController)
			studentEnrollments.PATCH("/:id", mod.StudentEnrollments.Ctl.StudentEnrollmentsUpdate)
			studentEnrollments.DELETE("/:id", mod.StudentEnrollments.Ctl.StudentEnrollmentsDelete)
		}

		enrollmentStatusHistories := enrollments.Group("/enrollment-status-histories")
		{
			enrollmentStatusHistories.GET("", mod.EnrollmentStatusHistories.Ctl.EnrollmentStatusHistoriesList)
			enrollmentStatusHistories.GET("/:id", mod.EnrollmentStatusHistories.Ctl.EnrollmentStatusHistoriesInfo)
			enrollmentStatusHistories.POST("", mod.EnrollmentStatusHistories.Ctl.CreateEnrollmentStatusHistoryController)
			enrollmentStatusHistories.PATCH("/:id", mod.EnrollmentStatusHistories.Ctl.EnrollmentStatusHistoriesUpdate)
			enrollmentStatusHistories.DELETE("/:id", mod.EnrollmentStatusHistories.Ctl.EnrollmentStatusHistoriesDelete)
		}

		enrollmentSubjects := enrollments.Group("/enrollment-subjects")
		{
			enrollmentSubjects.GET("", mod.EnrollmentSubjects.Ctl.EnrollmentSubjectsList)
			enrollmentSubjects.GET("/:id", mod.EnrollmentSubjects.Ctl.EnrollmentSubjectsInfo)
			enrollmentSubjects.POST("", mod.EnrollmentSubjects.Ctl.CreateEnrollmentSubjectController)
			enrollmentSubjects.PATCH("/:id", mod.EnrollmentSubjects.Ctl.EnrollmentSubjectsUpdate)
			enrollmentSubjects.DELETE("/:id", mod.EnrollmentSubjects.Ctl.EnrollmentSubjectsDelete)
		}
	}
}

func apiAttendance(r *gin.RouterGroup, mod *modules.Modules) {
	// Attendance routes (authentication required) e.g. for student attendance, etc.
	attendance := r.Group("/attendance")
	{
		sessions := attendance.Group("/sessions")
		{
			sessions.GET("", mod.AttendanceSessions.Ctl.AttendanceSessionsList)
			sessions.GET("/:id", mod.AttendanceSessions.Ctl.AttendanceSessionsInfo)
			sessions.GET("/:id/roster", mod.AttendanceSessions.Ctl.AttendanceSessionsRoster)
			sessions.POST("", mod.AttendanceSessions.Ctl.CreateAttendanceSessionController)
			sessions.PATCH("/:id", mod.AttendanceSessions.Ctl.AttendanceSessionsUpdate)
			sessions.DELETE("/:id", mod.AttendanceSessions.Ctl.AttendanceSessionsDelete)
		}

		records := attendance.Group("/records")
		{
			records.GET("", mod.AttendanceRecords.Ctl.AttendanceRecordsList)
			records.GET("/:id", mod.AttendanceRecords.Ctl.AttendanceRecordsInfo)
			records.POST("", mod.AttendanceRecords.Ctl.CreateAttendanceRecordController)
			records.PATCH("/:id", mod.AttendanceRecords.Ctl.AttendanceRecordsUpdate)
			records.PATCH("/bulk-mark", mod.AttendanceRecords.Ctl.AttendanceRecordsBulkMark)
			records.DELETE("/:id", mod.AttendanceRecords.Ctl.AttendanceRecordsDelete)
		}

		recordLogs := attendance.Group("/record-logs")
		{
			recordLogs.GET("", mod.AttendanceRecordLogs.Ctl.AttendanceRecordLogsList)
			recordLogs.GET("/:id", mod.AttendanceRecordLogs.Ctl.AttendanceRecordLogsInfo)
			recordLogs.POST("", mod.AttendanceRecordLogs.Ctl.CreateAttendanceRecordLogController)
			recordLogs.PATCH("/:id", mod.AttendanceRecordLogs.Ctl.AttendanceRecordLogsUpdate)
			recordLogs.DELETE("/:id", mod.AttendanceRecordLogs.Ctl.AttendanceRecordLogsDelete)
		}
	}
}

func apiApproval(r *gin.RouterGroup, mod *modules.Modules) {
	approval := r.Group("/approval")
	{
		requests := approval.Group("/requests")
		{
			requests.GET("", mod.Approvals.Ctl.ApprovalRequestsList)
			requests.GET("/:id", mod.Approvals.Ctl.ApprovalRequestsInfo)
			requests.GET("/:id/actions", mod.Approvals.Ctl.ApprovalRequestsActionsList)
			requests.POST("", mod.Approvals.Ctl.CreateApprovalRequestController)
			requests.POST("/:id/submit", mod.Approvals.Ctl.ApprovalRequestsSubmit)
			requests.POST("/:id/approve", mod.Approvals.Ctl.ApprovalRequestsApprove)
			requests.POST("/:id/reject", mod.Approvals.Ctl.ApprovalRequestsReject)
			requests.POST("/:id/cancel", mod.Approvals.Ctl.ApprovalRequestsCancel)
		}
	}
}

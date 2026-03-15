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

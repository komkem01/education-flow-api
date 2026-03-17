package studentregistrationcases

import (
	"strconv"

	"eduflow/app/modules/auth"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type byIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type listRequest struct {
	base.RequestPaginate
	SchoolID         string `form:"school_id"`
	StudentID        string `form:"student_id"`
	Status           string `form:"status"`
	RegistrationType string `form:"registration_type"`
}

type actionRequest struct {
	IdempotencyKey *string        `json:"idempotency_key"`
	Comment        *string        `json:"comment"`
	Metadata       map[string]any `json:"metadata"`
}

type createUpdateRequest struct {
	CaseNo            *string                   `json:"case_no"`
	SchoolID          *uuid.UUID                `json:"school_id"`
	StudentID         *uuid.UUID                `json:"student_id"`
	RegistrationType  *string                   `json:"registration_type"`
	EffectiveDate     *string                   `json:"effective_date"`
	Reason            *string                   `json:"reason"`
	Email             *string                   `json:"email"`
	Password          *string                   `json:"password"`
	Core              *coreRequest              `json:"core"`
	Addresses         []addressRequest          `json:"addresses"`
	Health            *healthRequest            `json:"health"`
	Guardians         []guardianRequest         `json:"guardians"`
	StudentGuardians  []studentGuardianRequest  `json:"student_guardians"`
	PreviousEducation *previousEducationRequest `json:"previous_education"`
	FamilyEconomic    *familyEconomicRequest    `json:"family_economic"`
	Documents         []documentRequest         `json:"documents"`
}

type coreRequest struct {
	MemberID         *uuid.UUID `json:"member_id"`
	StudentID        *uuid.UUID `json:"student_id"`
	GenderID         uuid.UUID  `json:"gender_id"`
	PrefixID         uuid.UUID  `json:"prefix_id"`
	AdvisorTeacherID *uuid.UUID `json:"advisor_teacher_id"`
	FirstNameTH      string     `json:"first_name_th"`
	LastNameTH       string     `json:"last_name_th"`
	FirstNameEN      *string    `json:"first_name_en"`
	LastNameEN       *string    `json:"last_name_en"`
	CitizenID        *string    `json:"citizen_id"`
	Phone            *string    `json:"phone"`
	IsActiveTarget   *bool      `json:"is_active_target"`
}

type addressRequest struct {
	AddressType string  `json:"address_type"`
	HouseNo     string  `json:"house_no"`
	Village     *string `json:"village"`
	Road        *string `json:"road"`
	Province    string  `json:"province"`
	District    string  `json:"district"`
	Subdistrict string  `json:"subdistrict"`
	PostalCode  string  `json:"postal_code"`
	Country     *string `json:"country"`
	IsPrimary   *bool   `json:"is_primary"`
	SortOrder   *int    `json:"sort_order"`
}

type healthRequest struct {
	BloodType            *string `json:"blood_type"`
	AllergyInfo          *string `json:"allergy_info"`
	ChronicDisease       *string `json:"chronic_disease"`
	MedicalNote          *string `json:"medical_note"`
	DisabilityFlag       *bool   `json:"disability_flag"`
	DisabilityDetail     *string `json:"disability_detail"`
	SpecialSupportFlag   *bool   `json:"special_support_flag"`
	SpecialSupportDetail *string `json:"special_support_detail"`
}

type guardianRequest struct {
	RowID            *uuid.UUID `json:"row_id"`
	GenderID         uuid.UUID  `json:"gender_id"`
	PrefixID         uuid.UUID  `json:"prefix_id"`
	FirstNameTH      string     `json:"first_name_th"`
	LastNameTH       string     `json:"last_name_th"`
	FirstNameEN      *string    `json:"first_name_en"`
	LastNameEN       *string    `json:"last_name_en"`
	CitizenID        *string    `json:"citizen_id"`
	Phone            string     `json:"phone"`
	Occupation       *string    `json:"occupation"`
	Employer         *string    `json:"employer"`
	MonthlyIncome    *float64   `json:"monthly_income"`
	AnnualIncome     *float64   `json:"annual_income"`
	EducationLevel   *string    `json:"education_level"`
	RelationshipText *string    `json:"relationship_text"`
	IsActiveTarget   *bool      `json:"is_active_target"`
	SortOrder        *int       `json:"sort_order"`
}

type studentGuardianRequest struct {
	GuardianRowID      uuid.UUID `json:"guardian_row_id"`
	Relationship       string    `json:"relationship"`
	IsMainGuardian     *bool     `json:"is_main_guardian"`
	CanPickup          *bool     `json:"can_pickup"`
	IsEmergencyContact *bool     `json:"is_emergency_contact"`
	Note               *string   `json:"note"`
	SortOrder          *int      `json:"sort_order"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	user, ok := auth.CurrentUserFromGin(ctx)
	if !ok || user.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var req createUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), user.Member.ID, user.Member.Role, mapUpsertInput(&req))
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req listRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var schoolID *uuid.UUID
	if req.SchoolID != "" {
		parsed, err := uuid.Parse(req.SchoolID)
		if err != nil {
			base.BadRequest(ctx, "invalid-school-id", nil)
			return
		}
		schoolID = &parsed
	}
	var studentID *uuid.UUID
	if req.StudentID != "" {
		parsed, err := uuid.Parse(req.StudentID)
		if err != nil {
			base.BadRequest(ctx, "invalid-student-id", nil)
			return
		}
		studentID = &parsed
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, schoolID, studentID, strPtr(req.Status), strPtr(req.RegistrationType))
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) Info(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, ok := parseIDFromURI(ctx)
	if !ok {
		return
	}

	item, err := c.svc.Info(ctx.Request.Context(), id)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-info-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) Update(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, ok := parseIDFromURI(ctx)
	if !ok {
		return
	}
	user, ok := auth.CurrentUserFromGin(ctx)
	if !ok || user.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var req createUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	item, err := c.svc.Update(ctx.Request.Context(), id, user.Member.Role, mapUpsertInput(&req))
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) Delete(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, ok := parseIDFromURI(ctx)
	if !ok {
		return
	}
	user, ok := auth.CurrentUserFromGin(ctx)
	if !ok || user.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	if err := c.svc.Delete(ctx.Request.Context(), id, user.Member.Role); err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-delete-failed")
		return
	}

	base.Success(ctx, gin.H{"id": id}, "success")
}

func (c *Controller) Submit(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, ok := parseIDFromURI(ctx)
	if !ok {
		return
	}
	user, ok := auth.CurrentUserFromGin(ctx)
	if !ok || user.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var req actionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	item, err := c.svc.Submit(ctx.Request.Context(), id, user.Member.ID, user.Member.Role, req.Comment)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-submit-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) Approve(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, ok := parseIDFromURI(ctx)
	if !ok {
		return
	}
	user, ok := auth.CurrentUserFromGin(ctx)
	if !ok || user.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var req actionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	item, err := c.svc.Approve(ctx.Request.Context(), id, user.Member.ID, user.Member.Role, req.Comment, req.IdempotencyKey, req.Metadata)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-approve-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) Reject(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, ok := parseIDFromURI(ctx)
	if !ok {
		return
	}
	user, ok := auth.CurrentUserFromGin(ctx)
	if !ok || user.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var req actionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	item, err := c.svc.Reject(ctx.Request.Context(), id, user.Member.ID, user.Member.Role, req.Comment)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-reject-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) Apply(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, ok := parseIDFromURI(ctx)
	if !ok {
		return
	}
	user, ok := auth.CurrentUserFromGin(ctx)
	if !ok || user.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var req actionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	item, err := c.svc.Apply(ctx.Request.Context(), id, user.Member.ID, user.Member.Role, req.Comment)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-apply-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func parseIDFromURI(ctx *gin.Context) (uuid.UUID, bool) {
	var req byIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return uuid.Nil, false
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return uuid.Nil, false
	}
	return id, true
}

func strPtr(v string) *string {
	v = string([]byte(v))
	if v == "" {
		return nil
	}
	return &v
}

func mapUpsertInput(req *createUpdateRequest) *UpsertInput {
	if req == nil {
		return nil
	}

	out := &UpsertInput{
		CaseNo:           req.CaseNo,
		SchoolID:         req.SchoolID,
		StudentID:        req.StudentID,
		RegistrationType: req.RegistrationType,
		EffectiveDate:    req.EffectiveDate,
		Reason:           req.Reason,
		Email:            req.Email,
		Password:         req.Password,
	}

	if req.Core != nil {
		out.Core = &CoreInput{
			MemberID:         req.Core.MemberID,
			StudentID:        req.Core.StudentID,
			GenderID:         req.Core.GenderID,
			PrefixID:         req.Core.PrefixID,
			AdvisorTeacherID: req.Core.AdvisorTeacherID,
			FirstNameTH:      req.Core.FirstNameTH,
			LastNameTH:       req.Core.LastNameTH,
			FirstNameEN:      req.Core.FirstNameEN,
			LastNameEN:       req.Core.LastNameEN,
			CitizenID:        req.Core.CitizenID,
			Phone:            req.Core.Phone,
			IsActiveTarget:   req.Core.IsActiveTarget,
		}
	}

	for _, row := range req.Addresses {
		out.Addresses = append(out.Addresses, AddressInput{
			AddressType: row.AddressType,
			HouseNo:     row.HouseNo,
			Village:     row.Village,
			Road:        row.Road,
			Province:    row.Province,
			District:    row.District,
			Subdistrict: row.Subdistrict,
			PostalCode:  row.PostalCode,
			Country:     row.Country,
			IsPrimary:   row.IsPrimary,
			SortOrder:   row.SortOrder,
		})
	}

	if req.Health != nil {
		out.Health = &HealthInput{
			BloodType:            req.Health.BloodType,
			AllergyInfo:          req.Health.AllergyInfo,
			ChronicDisease:       req.Health.ChronicDisease,
			MedicalNote:          req.Health.MedicalNote,
			DisabilityFlag:       req.Health.DisabilityFlag,
			DisabilityDetail:     req.Health.DisabilityDetail,
			SpecialSupportFlag:   req.Health.SpecialSupportFlag,
			SpecialSupportDetail: req.Health.SpecialSupportDetail,
		}
	}

	for _, row := range req.Guardians {
		out.Guardians = append(out.Guardians, GuardianInput{
			RowID:            row.RowID,
			GenderID:         row.GenderID,
			PrefixID:         row.PrefixID,
			FirstNameTH:      row.FirstNameTH,
			LastNameTH:       row.LastNameTH,
			FirstNameEN:      row.FirstNameEN,
			LastNameEN:       row.LastNameEN,
			CitizenID:        row.CitizenID,
			Phone:            row.Phone,
			Occupation:       row.Occupation,
			Employer:         row.Employer,
			MonthlyIncome:    row.MonthlyIncome,
			AnnualIncome:     row.AnnualIncome,
			EducationLevel:   row.EducationLevel,
			RelationshipText: row.RelationshipText,
			IsActiveTarget:   row.IsActiveTarget,
			SortOrder:        row.SortOrder,
		})
	}

	for _, row := range req.StudentGuardians {
		out.StudentGuardians = append(out.StudentGuardians, StudentGuardianInput{
			GuardianRowID:      row.GuardianRowID,
			Relationship:       row.Relationship,
			IsMainGuardian:     row.IsMainGuardian,
			CanPickup:          row.CanPickup,
			IsEmergencyContact: row.IsEmergencyContact,
			Note:               row.Note,
			SortOrder:          row.SortOrder,
		})
	}

	if req.PreviousEducation != nil {
		out.PreviousEducation = &PreviousEducationInput{
			PreviousSchoolName:     req.PreviousEducation.PreviousSchoolName,
			PreviousSchoolProvince: req.PreviousEducation.PreviousSchoolProvince,
			PreviousGradeLevel:     req.PreviousEducation.PreviousGradeLevel,
			GPA:                    req.PreviousEducation.GPA,
			TransferCertificateNo:  req.PreviousEducation.TransferCertificateNo,
			TransferDate:           req.PreviousEducation.TransferDate,
			TranscriptReceived:     req.PreviousEducation.TranscriptReceived,
			Remarks:                req.PreviousEducation.Remarks,
		}
	}

	if req.FamilyEconomic != nil {
		out.FamilyEconomic = &FamilyEconomicInput{
			HouseholdSize:          req.FamilyEconomic.HouseholdSize,
			HouseholdIncomeMonthly: req.FamilyEconomic.HouseholdIncomeMonthly,
			IncomeBracket:          req.FamilyEconomic.IncomeBracket,
			ScholarshipFlag:        req.FamilyEconomic.ScholarshipFlag,
			ScholarshipType:        req.FamilyEconomic.ScholarshipType,
			WelfareFlag:            req.FamilyEconomic.WelfareFlag,
			WelfareType:            req.FamilyEconomic.WelfareType,
			DebtFlag:               req.FamilyEconomic.DebtFlag,
			DebtDetail:             req.FamilyEconomic.DebtDetail,
		}
	}

	for _, row := range req.Documents {
		out.Documents = append(out.Documents, DocumentInput{
			DocType:        row.DocType,
			FileDocumentID: row.FileDocumentID,
			FileName:       row.FileName,
			MimeType:       row.MimeType,
			FileSizeBytes:  row.FileSizeBytes,
			IsRequired:     row.IsRequired,
			IsVerified:     row.IsVerified,
			VerifiedBy:     row.VerifiedBy,
			VerifiedAt:     row.VerifiedAt,
			Note:           row.Note,
		})
	}

	return out
}

func (c *Controller) StudentRegistrationCasesCreate(ctx *gin.Context)  { c.Create(ctx) }
func (c *Controller) StudentRegistrationCasesList(ctx *gin.Context)    { c.List(ctx) }
func (c *Controller) StudentRegistrationCasesInfo(ctx *gin.Context)    { c.Info(ctx) }
func (c *Controller) StudentRegistrationCasesUpdate(ctx *gin.Context)  { c.Update(ctx) }
func (c *Controller) StudentRegistrationCasesDelete(ctx *gin.Context)  { c.Delete(ctx) }
func (c *Controller) StudentRegistrationCasesSubmit(ctx *gin.Context)  { c.Submit(ctx) }
func (c *Controller) StudentRegistrationCasesApprove(ctx *gin.Context) { c.Approve(ctx) }
func (c *Controller) StudentRegistrationCasesReject(ctx *gin.Context)  { c.Reject(ctx) }
func (c *Controller) StudentRegistrationCasesApply(ctx *gin.Context)   { c.Apply(ctx) }

func boolFromQuery(ctx *gin.Context, key string) *bool {
	v := ctx.Query(key)
	if v == "" {
		return nil
	}
	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return nil
	}
	return &parsed
}

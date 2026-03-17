package studentregistrationcases

import (
	"eduflow/app/modules/auth"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type previousEducationRequest struct {
	PreviousSchoolName     *string  `json:"previous_school_name"`
	PreviousSchoolProvince *string  `json:"previous_school_province"`
	PreviousGradeLevel     *string  `json:"previous_grade_level"`
	GPA                    *float64 `json:"gpa"`
	TransferCertificateNo  *string  `json:"transfer_certificate_no"`
	TransferDate           *string  `json:"transfer_date"`
	TranscriptReceived     *bool    `json:"transcript_received"`
	Remarks                *string  `json:"remarks"`
}

type familyEconomicRequest struct {
	HouseholdSize          *int     `json:"household_size"`
	HouseholdIncomeMonthly *float64 `json:"household_income_monthly"`
	IncomeBracket          *string  `json:"income_bracket"`
	ScholarshipFlag        *bool    `json:"scholarship_flag"`
	ScholarshipType        *string  `json:"scholarship_type"`
	WelfareFlag            *bool    `json:"welfare_flag"`
	WelfareType            *string  `json:"welfare_type"`
	DebtFlag               *bool    `json:"debt_flag"`
	DebtDetail             *string  `json:"debt_detail"`
}

type documentRequest struct {
	DocType        string     `json:"doc_type" binding:"required"`
	FileDocumentID *uuid.UUID `json:"file_document_id"`
	FileName       *string    `json:"file_name"`
	MimeType       *string    `json:"mime_type"`
	FileSizeBytes  *int64     `json:"file_size_bytes"`
	IsRequired     *bool      `json:"is_required"`
	IsVerified     *bool      `json:"is_verified"`
	VerifiedBy     *uuid.UUID `json:"verified_by"`
	VerifiedAt     *string    `json:"verified_at"`
	Note           *string    `json:"note"`
}

type replaceDocumentsRequest struct {
	Documents []documentRequest `json:"documents"`
}

type ruleRequest struct {
	SchoolID          uuid.UUID `json:"school_id"`
	RegistrationType  string    `json:"registration_type"`
	FieldCode         string    `json:"field_code"`
	IsRequired        *bool     `json:"is_required"`
	ValidationRegex   *string   `json:"validation_regex"`
	ValidationMessage *string   `json:"validation_message"`
	ActiveFrom        *string   `json:"active_from"`
	ActiveTo          *string   `json:"active_to"`
}

type ruleListRequest struct {
	base.RequestPaginate
	SchoolID         string `form:"school_id"`
	RegistrationType string `form:"registration_type"`
}

func (c *Controller) GetPreviousEducation(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, ok := parseIDFromURI(ctx)
	if !ok {
		return
	}
	item, err := c.svc.GetPreviousEducation(ctx.Request.Context(), id)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-previous-education-get-failed")
		return
	}
	base.Success(ctx, item, "success")
}

func (c *Controller) UpsertPreviousEducation(ctx *gin.Context) {
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

	var req previousEducationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	item, err := c.svc.UpsertPreviousEducation(ctx.Request.Context(), id, user.Member.Role, &PreviousEducationInput{
		PreviousSchoolName:     req.PreviousSchoolName,
		PreviousSchoolProvince: req.PreviousSchoolProvince,
		PreviousGradeLevel:     req.PreviousGradeLevel,
		GPA:                    req.GPA,
		TransferCertificateNo:  req.TransferCertificateNo,
		TransferDate:           req.TransferDate,
		TranscriptReceived:     req.TranscriptReceived,
		Remarks:                req.Remarks,
	})
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-previous-education-upsert-failed")
		return
	}
	base.Success(ctx, item, "success")
}

func (c *Controller) GetFamilyEconomic(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, ok := parseIDFromURI(ctx)
	if !ok {
		return
	}
	item, err := c.svc.GetFamilyEconomic(ctx.Request.Context(), id)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-family-economic-get-failed")
		return
	}
	base.Success(ctx, item, "success")
}

func (c *Controller) UpsertFamilyEconomic(ctx *gin.Context) {
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

	var req familyEconomicRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	item, err := c.svc.UpsertFamilyEconomic(ctx.Request.Context(), id, user.Member.Role, &FamilyEconomicInput{
		HouseholdSize:          req.HouseholdSize,
		HouseholdIncomeMonthly: req.HouseholdIncomeMonthly,
		IncomeBracket:          req.IncomeBracket,
		ScholarshipFlag:        req.ScholarshipFlag,
		ScholarshipType:        req.ScholarshipType,
		WelfareFlag:            req.WelfareFlag,
		WelfareType:            req.WelfareType,
		DebtFlag:               req.DebtFlag,
		DebtDetail:             req.DebtDetail,
	})
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-family-economic-upsert-failed")
		return
	}
	base.Success(ctx, item, "success")
}

func (c *Controller) ListDocuments(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, ok := parseIDFromURI(ctx)
	if !ok {
		return
	}
	rows, err := c.svc.ListDocuments(ctx.Request.Context(), id)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-documents-list-failed")
		return
	}
	base.Success(ctx, rows, "success")
}

func (c *Controller) ReplaceDocuments(ctx *gin.Context) {
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

	var req replaceDocumentsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	rows := make([]DocumentInput, 0, len(req.Documents))
	for _, row := range req.Documents {
		rows = append(rows, DocumentInput{
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

	item, err := c.svc.ReplaceDocuments(ctx.Request.Context(), id, user.Member.Role, rows)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-case-documents-replace-failed")
		return
	}
	base.Success(ctx, item, "success")
}

func (c *Controller) ListRules(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req ruleListRequest
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

	items, page, err := c.svc.ListRules(ctx.Request.Context(), &req.RequestPaginate, schoolID, strPtr(req.RegistrationType))
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-rules-list-failed")
		return
	}
	base.Paginate(ctx, items, page)
}

func (c *Controller) GetRule(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	id, ok := parseIDFromURI(ctx)
	if !ok {
		return
	}
	item, err := c.svc.GetRule(ctx.Request.Context(), id)
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-rule-get-failed")
		return
	}
	base.Success(ctx, item, "success")
}

func (c *Controller) CreateRule(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	user, ok := auth.CurrentUserFromGin(ctx)
	if !ok || user.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var req ruleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}
	item, err := c.svc.CreateRule(ctx.Request.Context(), user.Member.Role, &RuleUpsertInput{
		SchoolID:          req.SchoolID,
		RegistrationType:  req.RegistrationType,
		FieldCode:         req.FieldCode,
		IsRequired:        req.IsRequired,
		ValidationRegex:   req.ValidationRegex,
		ValidationMessage: req.ValidationMessage,
		ActiveFrom:        req.ActiveFrom,
		ActiveTo:          req.ActiveTo,
	})
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-rule-create-failed")
		return
	}
	base.Success(ctx, item, "success")
}

func (c *Controller) UpdateRule(ctx *gin.Context) {
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

	var req ruleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	item, err := c.svc.UpdateRule(ctx.Request.Context(), id, user.Member.Role, &RuleUpsertInput{
		SchoolID:          req.SchoolID,
		RegistrationType:  req.RegistrationType,
		FieldCode:         req.FieldCode,
		IsRequired:        req.IsRequired,
		ValidationRegex:   req.ValidationRegex,
		ValidationMessage: req.ValidationMessage,
		ActiveFrom:        req.ActiveFrom,
		ActiveTo:          req.ActiveTo,
	})
	if err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-rule-update-failed")
		return
	}
	base.Success(ctx, item, "success")
}

func (c *Controller) DeleteRule(ctx *gin.Context) {
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

	if err := c.svc.DeleteRule(ctx.Request.Context(), id, user.Member.Role); err != nil {
		c.handleServiceError(ctx, log, err, "student-registration-rule-delete-failed")
		return
	}
	base.Success(ctx, gin.H{"id": id}, "success")
}

func (c *Controller) StudentRegistrationCasesPreviousEducationInfo(ctx *gin.Context) {
	c.GetPreviousEducation(ctx)
}
func (c *Controller) StudentRegistrationCasesPreviousEducationUpsert(ctx *gin.Context) {
	c.UpsertPreviousEducation(ctx)
}
func (c *Controller) StudentRegistrationCasesFamilyEconomicInfo(ctx *gin.Context) {
	c.GetFamilyEconomic(ctx)
}
func (c *Controller) StudentRegistrationCasesFamilyEconomicUpsert(ctx *gin.Context) {
	c.UpsertFamilyEconomic(ctx)
}
func (c *Controller) StudentRegistrationCasesDocumentsList(ctx *gin.Context) { c.ListDocuments(ctx) }
func (c *Controller) StudentRegistrationCasesDocumentsReplace(ctx *gin.Context) {
	c.ReplaceDocuments(ctx)
}
func (c *Controller) StudentRegistrationRulesList(ctx *gin.Context)   { c.ListRules(ctx) }
func (c *Controller) StudentRegistrationRulesInfo(ctx *gin.Context)   { c.GetRule(ctx) }
func (c *Controller) StudentRegistrationRulesCreate(ctx *gin.Context) { c.CreateRule(ctx) }
func (c *Controller) StudentRegistrationRulesUpdate(ctx *gin.Context) { c.UpdateRule(ctx) }
func (c *Controller) StudentRegistrationRulesDelete(ctx *gin.Context) { c.DeleteRule(ctx) }

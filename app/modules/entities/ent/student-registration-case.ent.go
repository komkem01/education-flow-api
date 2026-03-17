package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StudentRegistrationType string

const (
	StudentRegistrationTypeNewEnrollment StudentRegistrationType = "new_enrollment"
	StudentRegistrationTypeTransferIn    StudentRegistrationType = "transfer_in"
	StudentRegistrationTypeTransferOut   StudentRegistrationType = "transfer_out"
	StudentRegistrationTypeLeaveAbsence  StudentRegistrationType = "leave_of_absence"
	StudentRegistrationTypeWithdrawal    StudentRegistrationType = "withdrawal"
	StudentRegistrationTypeReEnrollment  StudentRegistrationType = "re_enrollment"
)

type StudentRegistrationCaseStatus string

const (
	StudentRegistrationCaseStatusDraft     StudentRegistrationCaseStatus = "draft"
	StudentRegistrationCaseStatusPending   StudentRegistrationCaseStatus = "pending"
	StudentRegistrationCaseStatusApproved  StudentRegistrationCaseStatus = "approved"
	StudentRegistrationCaseStatusRejected  StudentRegistrationCaseStatus = "rejected"
	StudentRegistrationCaseStatusCancelled StudentRegistrationCaseStatus = "cancelled"
	StudentRegistrationCaseStatusApplied   StudentRegistrationCaseStatus = "applied"
)

type RegistrationAddressType string

const (
	RegistrationAddressTypeCurrent    RegistrationAddressType = "current"
	RegistrationAddressTypeRegistered RegistrationAddressType = "registered"
	RegistrationAddressTypeContact    RegistrationAddressType = "contact"
)

type RegistrationDocumentType string

const (
	RegistrationDocumentTypeTranscript       RegistrationDocumentType = "transcript"
	RegistrationDocumentTypeTransferLetter   RegistrationDocumentType = "transfer_letter"
	RegistrationDocumentTypeHouseholdReg     RegistrationDocumentType = "household_register"
	RegistrationDocumentTypeBirthCertificate RegistrationDocumentType = "birth_certificate"
	RegistrationDocumentTypeIDCard           RegistrationDocumentType = "id_card"
	RegistrationDocumentTypeMedicalCert      RegistrationDocumentType = "medical_certificate"
	RegistrationDocumentTypePhoto            RegistrationDocumentType = "photo"
	RegistrationDocumentTypeOther            RegistrationDocumentType = "other"
)

type RegistrationIncomeBracket string

const (
	RegistrationIncomeBracketUnder5000   RegistrationIncomeBracket = "under_5000"
	RegistrationIncomeBracket5001_10000  RegistrationIncomeBracket = "5001_10000"
	RegistrationIncomeBracket10001_20000 RegistrationIncomeBracket = "10001_20000"
	RegistrationIncomeBracket20001_40000 RegistrationIncomeBracket = "20001_40000"
	RegistrationIncomeBracket40001_60000 RegistrationIncomeBracket = "40001_60000"
	RegistrationIncomeBracketAbove60000  RegistrationIncomeBracket = "above_60000"
)

type StudentRegistrationCase struct {
	bun.BaseModel `bun:"table:student_registration_cases,alias:src"`

	ID               uuid.UUID                     `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	CaseNo           string                        `bun:"case_no,notnull,type:varchar(50)"`
	SchoolID         uuid.UUID                     `bun:"school_id,notnull,type:uuid"`
	StudentID        *uuid.UUID                    `bun:"student_id,type:uuid"`
	RegistrationType StudentRegistrationType       `bun:"registration_type,notnull,type:student_registration_type"`
	Status           StudentRegistrationCaseStatus `bun:"status,notnull,type:student_registration_case_status,default:draft"`
	RequestedBy      uuid.UUID                     `bun:"requested_by,notnull,type:uuid"`
	RequestedByRole  ApprovalActorRole             `bun:"requested_by_role,notnull,type:approval_actor_role"`
	ApprovedBy       *uuid.UUID                    `bun:"approved_by,type:uuid"`
	RejectedBy       *uuid.UUID                    `bun:"rejected_by,type:uuid"`
	RequestedAt      time.Time                     `bun:"requested_at,notnull,default:current_timestamp"`
	SubmittedAt      *time.Time                    `bun:"submitted_at"`
	ApprovedAt       *time.Time                    `bun:"approved_at"`
	RejectedAt       *time.Time                    `bun:"rejected_at"`
	EffectiveDate    *time.Time                    `bun:"effective_date,type:date"`
	Reason           *string                       `bun:"reason"`
	RejectionReason  *string                       `bun:"rejection_reason"`
	CreatedAt        time.Time                     `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt        time.Time                     `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt        *time.Time                    `bun:"deleted_at,soft_delete"`
}

type StudentRegistrationCaseUpdate struct {
	CaseNo           *string
	SchoolID         *uuid.UUID
	StudentID        *uuid.UUID
	RegistrationType *StudentRegistrationType
	Status           *StudentRegistrationCaseStatus
	ApprovedBy       *uuid.UUID
	RejectedBy       *uuid.UUID
	SubmittedAt      *time.Time
	ApprovedAt       *time.Time
	RejectedAt       *time.Time
	EffectiveDate    *time.Time
	Reason           *string
	RejectionReason  *string
}

type StudentRegistrationStudentCore struct {
	bun.BaseModel `bun:"table:student_registration_student_core,alias:srsc"`

	ID                        uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	CaseID                    uuid.UUID  `bun:"case_id,notnull,type:uuid"`
	MemberID                  *uuid.UUID `bun:"member_id,type:uuid"`
	StudentID                 *uuid.UUID `bun:"student_id,type:uuid"`
	GenderID                  uuid.UUID  `bun:"gender_id,notnull,type:uuid"`
	PrefixID                  uuid.UUID  `bun:"prefix_id,notnull,type:uuid"`
	AdvisorTeacherID          *uuid.UUID `bun:"advisor_teacher_id,type:uuid"`
	FirstNameTH               string     `bun:"first_name_th,notnull,type:varchar(255)"`
	LastNameTH                string     `bun:"last_name_th,notnull,type:varchar(255)"`
	FirstNameEN               *string    `bun:"first_name_en,type:varchar(255)"`
	LastNameEN                *string    `bun:"last_name_en,type:varchar(255)"`
	CitizenID                 *string    `bun:"citizen_id,type:varchar(13)"`
	Phone                     *string    `bun:"phone,type:varchar(20)"`
	PendingMemberEmail        *string    `bun:"pending_member_email,type:text" json:"-"`
	PendingMemberPasswordHash *string    `bun:"pending_member_password_hash,type:text" json:"-"`
	IsActiveTarget            bool       `bun:"is_active_target,notnull,default:true"`
	CreatedAt                 time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt                 time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
}

type StudentRegistrationAddress struct {
	bun.BaseModel `bun:"table:student_registration_addresses,alias:sra"`

	ID          uuid.UUID               `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	CaseID      uuid.UUID               `bun:"case_id,notnull,type:uuid"`
	AddressType RegistrationAddressType `bun:"address_type,notnull,type:registration_address_type"`
	HouseNo     string                  `bun:"house_no,notnull"`
	Village     *string                 `bun:"village"`
	Road        *string                 `bun:"road"`
	Province    string                  `bun:"province,notnull"`
	District    string                  `bun:"district,notnull"`
	Subdistrict string                  `bun:"subdistrict,notnull"`
	PostalCode  string                  `bun:"postal_code,notnull,type:varchar(10)"`
	Country     string                  `bun:"country,notnull,type:varchar(10)"`
	IsPrimary   bool                    `bun:"is_primary,notnull,default:false"`
	SortOrder   int                     `bun:"sort_order,notnull,default:1"`
	CreatedAt   time.Time               `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   time.Time               `bun:"updated_at,notnull,default:current_timestamp"`
}

type StudentRegistrationHealth struct {
	bun.BaseModel `bun:"table:student_registration_health,alias:srh"`

	ID                   uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	CaseID               uuid.UUID `bun:"case_id,notnull,type:uuid"`
	BloodType            *string   `bun:"blood_type,type:blood_type"`
	AllergyInfo          *string   `bun:"allergy_info"`
	ChronicDisease       *string   `bun:"chronic_disease"`
	MedicalNote          *string   `bun:"medical_note"`
	DisabilityFlag       bool      `bun:"disability_flag,notnull,default:false"`
	DisabilityDetail     *string   `bun:"disability_detail"`
	SpecialSupportFlag   bool      `bun:"special_support_flag,notnull,default:false"`
	SpecialSupportDetail *string   `bun:"special_support_detail"`
	CreatedAt            time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt            time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type StudentRegistrationGuardian struct {
	bun.BaseModel `bun:"table:student_registration_guardians,alias:srg"`

	ID               uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	CaseID           uuid.UUID `bun:"case_id,notnull,type:uuid"`
	GenderID         uuid.UUID `bun:"gender_id,notnull,type:uuid"`
	PrefixID         uuid.UUID `bun:"prefix_id,notnull,type:uuid"`
	FirstNameTH      string    `bun:"first_name_th,notnull,type:varchar(255)"`
	LastNameTH       string    `bun:"last_name_th,notnull,type:varchar(255)"`
	FirstNameEN      *string   `bun:"first_name_en,type:varchar(255)"`
	LastNameEN       *string   `bun:"last_name_en,type:varchar(255)"`
	CitizenID        *string   `bun:"citizen_id,type:varchar(13)"`
	Phone            string    `bun:"phone,notnull,type:varchar(20)"`
	Occupation       *string   `bun:"occupation,type:varchar(255)"`
	Employer         *string   `bun:"employer,type:varchar(255)"`
	MonthlyIncome    *float64  `bun:"monthly_income,type:numeric(12,2)"`
	AnnualIncome     *float64  `bun:"annual_income,type:numeric(12,2)"`
	EducationLevel   *string   `bun:"education_level,type:varchar(100)"`
	RelationshipText *string   `bun:"relationship_text,type:varchar(120)"`
	IsActiveTarget   bool      `bun:"is_active_target,notnull,default:true"`
	SortOrder        int       `bun:"sort_order,notnull,default:1"`
	CreatedAt        time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt        time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type StudentRegistrationStudentGuardian struct {
	bun.BaseModel `bun:"table:student_registration_student_guardians,alias:srsg"`

	ID                 uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	CaseID             uuid.UUID `bun:"case_id,notnull,type:uuid"`
	GuardianRowID      uuid.UUID `bun:"guardian_row_id,notnull,type:uuid"`
	Relationship       string    `bun:"relationship,notnull,type:guardian_relationship"`
	IsMainGuardian     bool      `bun:"is_main_guardian,notnull,default:false"`
	CanPickup          bool      `bun:"can_pickup,notnull,default:true"`
	IsEmergencyContact bool      `bun:"is_emergency_contact,notnull,default:false"`
	Note               *string   `bun:"note"`
	SortOrder          int       `bun:"sort_order,notnull,default:1"`
	CreatedAt          time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt          time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type StudentRegistrationPreviousEducation struct {
	bun.BaseModel `bun:"table:student_registration_previous_education,alias:srpe"`

	ID                     uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	CaseID                 uuid.UUID  `bun:"case_id,notnull,type:uuid"`
	PreviousSchoolName     *string    `bun:"previous_school_name,type:varchar(255)"`
	PreviousSchoolProvince *string    `bun:"previous_school_province,type:varchar(100)"`
	PreviousGradeLevel     *string    `bun:"previous_grade_level,type:varchar(50)"`
	GPA                    *float64   `bun:"gpa,type:numeric(4,2)"`
	TransferCertificateNo  *string    `bun:"transfer_certificate_no,type:varchar(100)"`
	TransferDate           *time.Time `bun:"transfer_date,type:date"`
	TranscriptReceived     bool       `bun:"transcript_received,notnull,default:false"`
	Remarks                *string    `bun:"remarks"`
	CreatedAt              time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt              time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
}

type StudentRegistrationFamilyEconomic struct {
	bun.BaseModel `bun:"table:student_registration_family_economic,alias:srfe"`

	ID                     uuid.UUID                  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	CaseID                 uuid.UUID                  `bun:"case_id,notnull,type:uuid"`
	HouseholdSize          *int                       `bun:"household_size"`
	HouseholdIncomeMonthly *float64                   `bun:"household_income_monthly,type:numeric(12,2)"`
	IncomeBracket          *RegistrationIncomeBracket `bun:"income_bracket,type:registration_income_bracket"`
	ScholarshipFlag        bool                       `bun:"scholarship_flag,notnull,default:false"`
	ScholarshipType        *string                    `bun:"scholarship_type,type:varchar(120)"`
	WelfareFlag            bool                       `bun:"welfare_flag,notnull,default:false"`
	WelfareType            *string                    `bun:"welfare_type,type:varchar(120)"`
	DebtFlag               bool                       `bun:"debt_flag,notnull,default:false"`
	DebtDetail             *string                    `bun:"debt_detail"`
	CreatedAt              time.Time                  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt              time.Time                  `bun:"updated_at,notnull,default:current_timestamp"`
}

type StudentRegistrationDocument struct {
	bun.BaseModel `bun:"table:student_registration_documents,alias:srd"`

	ID             uuid.UUID                `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	CaseID         uuid.UUID                `bun:"case_id,notnull,type:uuid"`
	DocType        RegistrationDocumentType `bun:"doc_type,notnull,type:registration_document_type"`
	FileDocumentID *uuid.UUID               `bun:"file_document_id,type:uuid"`
	FileName       *string                  `bun:"file_name"`
	MimeType       *string                  `bun:"mime_type"`
	FileSizeBytes  *int64                   `bun:"file_size_bytes"`
	IsRequired     bool                     `bun:"is_required,notnull,default:false"`
	IsVerified     bool                     `bun:"is_verified,notnull,default:false"`
	VerifiedBy     *uuid.UUID               `bun:"verified_by,type:uuid"`
	VerifiedAt     *time.Time               `bun:"verified_at"`
	Note           *string                  `bun:"note"`
	CreatedAt      time.Time                `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt      time.Time                `bun:"updated_at,notnull,default:current_timestamp"`
}

type StudentRegistrationRule struct {
	bun.BaseModel `bun:"table:student_registration_rules,alias:srr"`

	ID                uuid.UUID               `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID          uuid.UUID               `bun:"school_id,notnull,type:uuid"`
	RegistrationType  StudentRegistrationType `bun:"registration_type,notnull,type:student_registration_type"`
	FieldCode         string                  `bun:"field_code,notnull,type:varchar(120)"`
	IsRequired        bool                    `bun:"is_required,notnull,default:false"`
	ValidationRegex   *string                 `bun:"validation_regex"`
	ValidationMessage *string                 `bun:"validation_message"`
	ActiveFrom        time.Time               `bun:"active_from,notnull,type:date"`
	ActiveTo          *time.Time              `bun:"active_to,type:date"`
	CreatedAt         time.Time               `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt         time.Time               `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt         *time.Time              `bun:"deleted_at,soft_delete"`
}

type StudentRegistrationRuleUpdate struct {
	SchoolID          *uuid.UUID
	RegistrationType  *StudentRegistrationType
	FieldCode         *string
	IsRequired        *bool
	ValidationRegex   *string
	ValidationMessage *string
	ActiveFrom        *time.Time
	ActiveTo          *time.Time
}

type StudentRegistrationAuditLog struct {
	bun.BaseModel `bun:"table:student_registration_audit_logs,alias:sral"`

	ID        uuid.UUID                      `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	CaseID    uuid.UUID                      `bun:"case_id,notnull,type:uuid"`
	Action    string                         `bun:"action,notnull,type:varchar(64)"`
	ActorID   uuid.UUID                      `bun:"actor_id,notnull,type:uuid"`
	ActorRole *ApprovalActorRole             `bun:"actor_role,type:approval_actor_role"`
	OldStatus *StudentRegistrationCaseStatus `bun:"old_status,type:student_registration_case_status"`
	NewStatus *StudentRegistrationCaseStatus `bun:"new_status,type:student_registration_case_status"`
	Comment   *string                        `bun:"comment"`
	CreatedAt time.Time                      `bun:"created_at,notnull,default:current_timestamp"`
}

type StudentRegistrationCaseBundle struct {
	Core              *StudentRegistrationStudentCore       `json:"core,omitempty"`
	Addresses         []*StudentRegistrationAddress         `json:"addresses,omitempty"`
	Health            *StudentRegistrationHealth            `json:"health,omitempty"`
	Guardians         []*StudentRegistrationGuardian        `json:"guardians,omitempty"`
	StudentGuardians  []*StudentRegistrationStudentGuardian `json:"student_guardians,omitempty"`
	PreviousEducation *StudentRegistrationPreviousEducation `json:"previous_education,omitempty"`
	FamilyEconomic    *StudentRegistrationFamilyEconomic    `json:"family_economic,omitempty"`
	Documents         []*StudentRegistrationDocument        `json:"documents,omitempty"`
}

type StudentRegistrationCaseDetail struct {
	Case              *StudentRegistrationCase
	Core              *StudentRegistrationStudentCore
	Addresses         []*StudentRegistrationAddress
	Health            *StudentRegistrationHealth
	Guardians         []*StudentRegistrationGuardian
	StudentGuardians  []*StudentRegistrationStudentGuardian
	PreviousEducation *StudentRegistrationPreviousEducation
	FamilyEconomic    *StudentRegistrationFamilyEconomic
	Documents         []*StudentRegistrationDocument
}

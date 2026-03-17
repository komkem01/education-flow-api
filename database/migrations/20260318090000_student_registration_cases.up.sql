SET statement_timeout = 0;

--bun:split

create type student_registration_type as enum (
    'new_enrollment',
    'transfer_in',
    'transfer_out',
    'leave_of_absence',
    'withdrawal',
    're_enrollment'
);

create type student_registration_case_status as enum (
    'draft',
    'pending',
    'approved',
    'rejected',
    'cancelled',
    'applied'
);

create type registration_address_type as enum (
    'current',
    'registered',
    'contact'
);

create type registration_document_type as enum (
    'transcript',
    'transfer_letter',
    'household_register',
    'birth_certificate',
    'id_card',
    'medical_certificate',
    'photo',
    'other'
);

create type registration_income_bracket as enum (
    'under_5000',
    '5001_10000',
    '10001_20000',
    '20001_40000',
    '40001_60000',
    'above_60000'
);

create table student_registration_cases (
    id uuid primary key default gen_random_uuid(),
    case_no varchar(50) not null,
    school_id uuid not null references schools(id),
    student_id uuid references member_students(id),
    registration_type student_registration_type not null,
    status student_registration_case_status not null default 'draft',
    requested_by uuid not null references members(id),
    requested_by_role approval_actor_role not null,
    approved_by uuid references members(id),
    rejected_by uuid references members(id),
    requested_at timestamptz not null default now(),
    submitted_at timestamptz,
    approved_at timestamptz,
    rejected_at timestamptz,
    effective_date date,
    reason text,
    rejection_reason text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (school_id, case_no)
);

create index idx_student_registration_cases_status on student_registration_cases(status);
create index idx_student_registration_cases_type on student_registration_cases(registration_type);
create index idx_student_registration_cases_school_id on student_registration_cases(school_id);
create index idx_student_registration_cases_student_id on student_registration_cases(student_id);
create index idx_student_registration_cases_requested_at on student_registration_cases(requested_at desc);

create table student_registration_student_core (
    id uuid primary key default gen_random_uuid(),
    case_id uuid not null unique references student_registration_cases(id) on delete cascade,
    member_id uuid references members(id),
    student_id uuid references member_students(id),
    gender_id uuid not null references genders(id),
    prefix_id uuid not null references prefixes(id),
    advisor_teacher_id uuid references member_teachers(id),
    first_name_th varchar(255) not null,
    last_name_th varchar(255) not null,
    first_name_en varchar(255),
    last_name_en varchar(255),
    citizen_id varchar(13),
    phone varchar(20),
    is_active_target boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    check ((first_name_en is null and last_name_en is null) or (first_name_en is not null and last_name_en is not null)),
    check (citizen_id is null or length(citizen_id) = 13)
);

create index idx_student_registration_student_core_case_id on student_registration_student_core(case_id);
create index idx_student_registration_student_core_student_id on student_registration_student_core(student_id);

create table student_registration_addresses (
    id uuid primary key default gen_random_uuid(),
    case_id uuid not null references student_registration_cases(id) on delete cascade,
    address_type registration_address_type not null,
    house_no text not null,
    village text,
    road text,
    province text not null,
    district text not null,
    subdistrict text not null,
    postal_code varchar(10) not null,
    country varchar(10) not null default 'TH',
    is_primary boolean not null default false,
    sort_order int not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    check (sort_order > 0)
);

create index idx_student_registration_addresses_case_id on student_registration_addresses(case_id);
create index idx_student_registration_addresses_case_type on student_registration_addresses(case_id, address_type);
create unique index ux_student_registration_addresses_primary_per_case_type
    on student_registration_addresses(case_id, address_type)
    where is_primary = true;

create table student_registration_health (
    id uuid primary key default gen_random_uuid(),
    case_id uuid not null unique references student_registration_cases(id) on delete cascade,
    blood_type blood_type,
    allergy_info text,
    chronic_disease text,
    medical_note text,
    disability_flag boolean not null default false,
    disability_detail text,
    special_support_flag boolean not null default false,
    special_support_detail text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table student_registration_previous_education (
    id uuid primary key default gen_random_uuid(),
    case_id uuid not null unique references student_registration_cases(id) on delete cascade,
    previous_school_name varchar(255),
    previous_school_province varchar(100),
    previous_grade_level varchar(50),
    gpa numeric(4,2),
    transfer_certificate_no varchar(100),
    transfer_date date,
    transcript_received boolean not null default false,
    remarks text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    check (gpa is null or (gpa >= 0 and gpa <= 4.00))
);

create table student_registration_guardians (
    id uuid primary key default gen_random_uuid(),
    case_id uuid not null references student_registration_cases(id) on delete cascade,
    gender_id uuid not null references genders(id),
    prefix_id uuid not null references prefixes(id),
    first_name_th varchar(255) not null,
    last_name_th varchar(255) not null,
    first_name_en varchar(255),
    last_name_en varchar(255),
    citizen_id varchar(13),
    phone varchar(20) not null,
    occupation varchar(255),
    employer varchar(255),
    monthly_income numeric(12,2),
    annual_income numeric(12,2),
    education_level varchar(100),
    relationship_text varchar(120),
    is_active_target boolean not null default true,
    sort_order int not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    check ((first_name_en is null and last_name_en is null) or (first_name_en is not null and last_name_en is not null)),
    check (citizen_id is null or length(citizen_id) = 13),
    check (monthly_income is null or monthly_income >= 0),
    check (annual_income is null or annual_income >= 0),
    check (sort_order > 0)
);

create index idx_student_registration_guardians_case_id on student_registration_guardians(case_id);

create table student_registration_student_guardians (
    id uuid primary key default gen_random_uuid(),
    case_id uuid not null references student_registration_cases(id) on delete cascade,
    guardian_row_id uuid not null references student_registration_guardians(id) on delete cascade,
    relationship guardian_relationship not null,
    is_main_guardian boolean not null default false,
    can_pickup boolean not null default true,
    is_emergency_contact boolean not null default false,
    note text,
    sort_order int not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    unique (case_id, guardian_row_id),
    check (sort_order > 0)
);

create index idx_student_registration_student_guardians_case_id on student_registration_student_guardians(case_id);
create unique index ux_student_registration_main_guardian_per_case
    on student_registration_student_guardians(case_id)
    where is_main_guardian = true;

create table student_registration_family_economic (
    id uuid primary key default gen_random_uuid(),
    case_id uuid not null unique references student_registration_cases(id) on delete cascade,
    household_size int,
    household_income_monthly numeric(12,2),
    income_bracket registration_income_bracket,
    scholarship_flag boolean not null default false,
    scholarship_type varchar(120),
    welfare_flag boolean not null default false,
    welfare_type varchar(120),
    debt_flag boolean not null default false,
    debt_detail text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    check (household_size is null or household_size > 0),
    check (household_income_monthly is null or household_income_monthly >= 0)
);

create table student_registration_documents (
    id uuid primary key default gen_random_uuid(),
    case_id uuid not null references student_registration_cases(id) on delete cascade,
    doc_type registration_document_type not null,
    file_document_id uuid references documents(id),
    file_name text,
    mime_type text,
    file_size_bytes bigint,
    is_required boolean not null default false,
    is_verified boolean not null default false,
    verified_by uuid references members(id),
    verified_at timestamptz,
    note text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    check (file_size_bytes is null or file_size_bytes >= 0)
);

create index idx_student_registration_documents_case_id on student_registration_documents(case_id);
create index idx_student_registration_documents_case_type on student_registration_documents(case_id, doc_type);

create table student_registration_rules (
    id uuid primary key default gen_random_uuid(),
    school_id uuid not null references schools(id),
    registration_type student_registration_type not null,
    field_code varchar(120) not null,
    is_required boolean not null default false,
    validation_regex text,
    validation_message text,
    active_from date not null default current_date,
    active_to date,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (school_id, registration_type, field_code, active_from)
);

create index idx_student_registration_rules_school_type on student_registration_rules(school_id, registration_type);

create table student_registration_audit_logs (
    id uuid primary key default gen_random_uuid(),
    case_id uuid not null references student_registration_cases(id) on delete cascade,
    action varchar(64) not null,
    actor_id uuid not null references members(id),
    actor_role approval_actor_role,
    old_status student_registration_case_status,
    new_status student_registration_case_status,
    comment text,
    created_at timestamptz not null default now()
);

create index idx_student_registration_audit_logs_case_id on student_registration_audit_logs(case_id);
create index idx_student_registration_audit_logs_created_at on student_registration_audit_logs(created_at desc);

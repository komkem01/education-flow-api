SET statement_timeout = 0;

--bun:split

create table subjects (
    id uuid primary key default gen_random_uuid(),
    school_id uuid not null references schools(id),
    subject_group_id uuid not null references subject_groups(id),
    code varchar(30) not null,
    name_th varchar(255) not null,
    name_en varchar(255),
    credit numeric(4,2) not null default 0,
    hours_per_week int,
    is_elective boolean not null default false,
    is_active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    constraint subjects_credit_non_negative check (credit >= 0),
    constraint subjects_hours_positive check (hours_per_week is null or hours_per_week > 0)
);

create unique index ux_subjects_school_code on subjects(school_id, code);
create unique index ux_subjects_school_group_name_th on subjects(school_id, subject_group_id, name_th);
create index idx_subjects_school_group on subjects(school_id, subject_group_id);

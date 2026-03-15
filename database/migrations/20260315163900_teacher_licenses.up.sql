SET statement_timeout = 0;

--bun:split

create type teacher_license_status as enum ('active', 'suspended', 'expired', 'revoked');

create table teacher_licenses (
    id uuid primary key default gen_random_uuid(),
    teacher_id uuid not null references member_teachers(id),
    license_no varchar(100) not null,
    issued_at date,
    expires_at date,
    license_status teacher_license_status not null default 'active',
    issued_by varchar(255),
    note text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    constraint teacher_licenses_date_range check (expires_at is null or issued_at is null or expires_at >= issued_at)
);

create unique index ux_teacher_licenses_license_no
    on teacher_licenses(license_no)
    where deleted_at is null;

create index idx_teacher_licenses_teacher_id
    on teacher_licenses(teacher_id);

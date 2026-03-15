SET statement_timeout = 0;

--bun:split

create type member_role as enum ('superadmin', 'admin', 'staff', 'teacher', 'student');

create table members (
    id uuid primary key default gen_random_uuid(),
    school_id uuid not null references schools(id),
    email text not null unique,
    password text not null,
    role member_role not null default 'admin',
    is_active boolean not null default false,
    last_login timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
)

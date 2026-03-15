SET statement_timeout = 0;

--bun:split

create table academic_years (
    id uuid primary key default gen_random_uuid(),
    school_id uuid not null references schools(id),
    year text not null,
    start_date date not null,
    end_date date not null,
    is_active boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
)

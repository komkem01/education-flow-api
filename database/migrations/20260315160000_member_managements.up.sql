SET statement_timeout = 0;

--bun:split

create table member_managements (
    id uuid primary key default gen_random_uuid(),
    member_id uuid not null unique references members(id),
    employee_code varchar(50) not null unique,
    position varchar(255) not null,
    start_work_date date not null,
    department_id uuid not null references departments(id),
    is_active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

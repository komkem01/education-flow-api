SET statement_timeout = 0;

--bun:split

create type teacher_request_type as enum (
    'edit',
    'delete',
    'other'
);
create type teacher_request_status as enum (
    'pending',
    'approved',
    'rejected'
);

create table teacher_requests (
    id uuid primary key default gen_random_uuid(),
    teacher_id uuid not null references member_teachers(id),
    request_type teacher_request_type not null,
    request_data jsonb not null,
    request_reason text,
    status teacher_request_status not null default 'pending',
    approved_by uuid references members(id),
    approved_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

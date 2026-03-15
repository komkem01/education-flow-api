SET statement_timeout = 0;

--bun:split

create type approval_actor_role as enum (
    'teacher',
    'admin'
);

create type approval_request_status as enum (
    'draft',
    'pending',
    'approved',
    'rejected',
    'cancelled'
);

create table approval_requests (
    id uuid primary key default gen_random_uuid(),
    request_type varchar(120) not null,
    subject_type varchar(120) not null,
    subject_id uuid,
    requested_by uuid not null references members(id),
    requested_by_role approval_actor_role not null,
    payload jsonb not null,
    current_status approval_request_status not null default 'draft',
    submitted_at timestamptz,
    resolved_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index idx_approval_requests_status_created_at on approval_requests(current_status, created_at desc);
create index idx_approval_requests_requested_by on approval_requests(requested_by);
create index idx_approval_requests_request_type on approval_requests(request_type);

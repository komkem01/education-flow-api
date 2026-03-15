SET statement_timeout = 0;

--bun:split

create table enrollment_status_histories (
    id uuid primary key default gen_random_uuid(),
    enrollment_id uuid not null references student_enrollments(id),
    from_status student_enrollment_status,
    to_status student_enrollment_status not null,
    changed_at timestamptz not null default now(),
    changed_by uuid references members(id),
    reason text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create index idx_enrollment_status_histories_enrollment_id on enrollment_status_histories(enrollment_id);
create index idx_enrollment_status_histories_changed_at on enrollment_status_histories(changed_at);

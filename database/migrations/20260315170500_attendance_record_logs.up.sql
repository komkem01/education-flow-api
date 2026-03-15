SET statement_timeout = 0;

--bun:split

create table attendance_record_logs (
    id uuid primary key default gen_random_uuid(),
    record_id uuid not null references attendance_records(id),
    old_status attendance_status,
    new_status attendance_status not null,
    changed_by uuid references members(id),
    changed_at timestamptz not null default now(),
    reason text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create index idx_attendance_record_logs_record_id on attendance_record_logs(record_id);
create index idx_attendance_record_logs_changed_at on attendance_record_logs(changed_at);

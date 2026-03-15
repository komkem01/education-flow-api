SET statement_timeout = 0;

--bun:split

create type attendance_status as enum ('present', 'late', 'absent', 'sick', 'leave', 'activity');
create type attendance_source as enum ('manual', 'qr', 'rfid', 'face', 'api');

create table attendance_records (
    id uuid primary key default gen_random_uuid(),
    session_id uuid not null references attendance_sessions(id),
    enrollment_id uuid not null references student_enrollments(id),
    status attendance_status not null default 'present',
    source attendance_source not null default 'manual',
    marked_at timestamptz not null default now(),
    remark text,
    marked_by uuid references members(id),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create unique index ux_attendance_records_session_enrollment on attendance_records(session_id, enrollment_id) where deleted_at is null;
create index idx_attendance_records_session on attendance_records(session_id);
create index idx_attendance_records_enrollment on attendance_records(enrollment_id);
create index idx_attendance_records_status on attendance_records(status);

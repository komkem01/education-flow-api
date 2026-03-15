SET statement_timeout = 0;

--bun:split

create type attendance_mode as enum ('homeroom', 'subject', 'activity');

create table attendance_sessions (
    id uuid primary key default gen_random_uuid(),
    school_id uuid not null references schools(id),
    academic_year_id uuid not null references academic_years(id),
    classroom_id uuid not null references classrooms(id),
    subject_id uuid references subjects(id),
    teacher_id uuid references member_teachers(id),
    session_date date not null,
    period_no int not null,
    mode attendance_mode not null default 'homeroom',
    started_at timestamptz not null default now(),
    closed_at timestamptz,
    note text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    constraint attendance_sessions_period_positive check (period_no > 0),
    constraint attendance_sessions_time_range check (closed_at is null or closed_at >= started_at)
);

create unique index ux_attendance_sessions_unique_slot on attendance_sessions(classroom_id, session_date, period_no, mode, coalesce(subject_id, '00000000-0000-0000-0000-000000000000'::uuid)) where deleted_at is null;
create index idx_attendance_sessions_school_year on attendance_sessions(school_id, academic_year_id);
create index idx_attendance_sessions_classroom_date on attendance_sessions(classroom_id, session_date);

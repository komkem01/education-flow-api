SET statement_timeout = 0;

--bun:split

create type student_enrollment_status as enum ('active', 'transferred', 'graduated', 'dropped');

create table student_enrollments (
    id uuid primary key default gen_random_uuid(),
    student_id uuid not null references member_students(id),
    school_id uuid not null references schools(id),
    academic_year_id uuid not null references academic_years(id),
    classroom_id uuid not null references classrooms(id),
    enrolled_at date,
    exited_at date,
    status student_enrollment_status not null default 'active',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    constraint student_enrollments_date_range check (exited_at is null or enrolled_at is null or exited_at >= enrolled_at)
);

create unique index ux_student_enrollments_active_per_student on student_enrollments(student_id) where status = 'active' and deleted_at is null;
create index idx_student_enrollments_student_id on student_enrollments(student_id);
create index idx_student_enrollments_school_year on student_enrollments(school_id, academic_year_id);
create index idx_student_enrollments_classroom_id on student_enrollments(classroom_id);

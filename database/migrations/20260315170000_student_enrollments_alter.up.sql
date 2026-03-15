SET statement_timeout = 0;

--bun:split

create type enrollment_type as enum ('new', 'transfer_in', 'repeat', 'return');
create type enrollment_exit_reason as enum ('transfer_out', 'graduated', 'dropped', 'leave');

alter table student_enrollments
    add column enrollment_type enrollment_type not null default 'new',
    add column exit_reason enrollment_exit_reason,
    add column exit_note text,
    add column previous_enrollment_id uuid references student_enrollments(id),
    add column roll_no varchar(10),
    add column approved_by uuid references members(id),
    add column approved_at timestamptz,
    add column approval_note text,
    add column created_by uuid references members(id),
    add column updated_by uuid references members(id);

alter table student_enrollments
    add constraint student_enrollments_status_exit_consistency check (
        status = 'active' or exited_at is not null
    );

create unique index ux_student_enrollments_student_academic_year on student_enrollments(student_id, academic_year_id) where deleted_at is null;
create index idx_student_enrollments_school_year_status on student_enrollments(school_id, academic_year_id, status);
create index idx_student_enrollments_classroom_status on student_enrollments(classroom_id, status);

SET statement_timeout = 0;

--bun:split

drop index if exists ux_student_enrollments_student_academic_year;
drop index if exists idx_student_enrollments_school_year_status;
drop index if exists idx_student_enrollments_classroom_status;

alter table student_enrollments
    drop constraint if exists student_enrollments_status_exit_consistency;

alter table student_enrollments
    drop column if exists enrollment_type,
    drop column if exists exit_reason,
    drop column if exists exit_note,
    drop column if exists previous_enrollment_id,
    drop column if exists roll_no,
    drop column if exists approved_by,
    drop column if exists approved_at,
    drop column if exists approval_note,
    drop column if exists created_by,
    drop column if exists updated_by;

drop type if exists enrollment_exit_reason;
drop type if exists enrollment_type;

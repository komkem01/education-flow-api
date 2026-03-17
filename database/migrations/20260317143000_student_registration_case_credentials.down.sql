SET statement_timeout = 0;

--bun:split

alter table student_registration_student_core
    drop column if exists pending_member_password_hash,
    drop column if exists pending_member_email;

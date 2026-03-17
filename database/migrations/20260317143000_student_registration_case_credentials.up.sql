SET statement_timeout = 0;

--bun:split

alter table student_registration_student_core
    add column if not exists pending_member_email text,
    add column if not exists pending_member_password_hash text;

SET statement_timeout = 0;

--bun:split

create type guardian_relationship as enum ('father', 'mother', 'guardian', 'other');

create table student_guardians (
    id uuid primary key default gen_random_uuid(),
    student_id uuid not null references member_students(id),
    guardian_id uuid not null references member_guardians(id),
    relationship guardian_relationship not null,
    is_main_guardian boolean not null default false,
    can_pickup boolean not null default true,
    is_emergency_contact boolean not null default false,
    note text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create unique index ux_student_guardians_student_guardian on student_guardians(student_id, guardian_id);
create unique index ux_student_guardians_main_guardian_per_student on student_guardians(student_id) where is_main_guardian = true and deleted_at is null;
create index idx_student_guardians_student_id on student_guardians(student_id);
create index idx_student_guardians_guardian_id on student_guardians(guardian_id);

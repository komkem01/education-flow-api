SET statement_timeout = 0;

--bun:split

create type blood_type as enum ('A', 'B', 'AB', 'O');

create table student_health_profiles (
    id uuid primary key default gen_random_uuid(),
    student_id uuid not null unique references member_students(id),
    blood_type blood_type,
    allergy_info text,
    chronic_disease text,
    medical_note text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create index idx_student_health_profiles_student_id on student_health_profiles(student_id);

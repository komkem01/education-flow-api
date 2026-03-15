SET statement_timeout = 0;

--bun:split

create table student_profiles (
    id uuid primary key default gen_random_uuid(),
    student_id uuid not null unique references member_students(id),
    birth_date date,
    nationality varchar(100),
    religion varchar(100),
    address_current text,
    address_registered text,
    emergency_contact_name varchar(255),
    emergency_contact_phone varchar(20),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create index idx_student_profiles_student_id on student_profiles(student_id);

SET statement_timeout = 0;

--bun:split

create type teacher_degree as enum ('ปริญญาตรี', 'ปริญญาโท', 'ปริญญาเอก');

create table teacher_educations (
    id uuid primary key default gen_random_uuid(),
    teacher_id uuid not null references member_teachers(id),
    degree teacher_degree not null,
    major varchar not null,
    university varchar not null,
    graduation_year varchar not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

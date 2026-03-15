SET statement_timeout = 0;

--bun:split

create table teacher_health_profiles (
    id uuid primary key default gen_random_uuid(),
    member_teacher_id uuid not null unique references member_teachers(id),
    blood_type blood_type,
    allergy_info text,
    chronic_disease text,
    medication_note text,
    fitness_for_work_note text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create index idx_teacher_health_profiles_member_teacher_id
    on teacher_health_profiles(member_teacher_id);

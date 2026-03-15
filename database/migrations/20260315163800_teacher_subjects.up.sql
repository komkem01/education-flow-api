SET statement_timeout = 0;

--bun:split

create type teacher_subject_role as enum ('primary', 'assistant');

create table teacher_subjects (
    id uuid primary key default gen_random_uuid(),
    teacher_id uuid not null references member_teachers(id),
    subject_id uuid not null references subjects(id),
    role teacher_subject_role not null default 'primary',
    is_active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create unique index ux_teacher_subjects_teacher_subject
    on teacher_subjects(teacher_id, subject_id)
    where deleted_at is null;

create index idx_teacher_subjects_subject_id
    on teacher_subjects(subject_id);

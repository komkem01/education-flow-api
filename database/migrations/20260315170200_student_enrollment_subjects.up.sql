SET statement_timeout = 0;

--bun:split

create table enrollment_subjects (
    id uuid primary key default gen_random_uuid(),
    enrollment_id uuid not null references student_enrollments(id),
    subject_id uuid not null references subjects(id),
    teacher_id uuid references member_teachers(id),
    is_primary boolean not null default false,
    status student_enrollment_status not null default 'active',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create unique index ux_enrollment_subjects_unique on enrollment_subjects(enrollment_id, subject_id) where deleted_at is null;
create index idx_enrollment_subjects_enrollment_id on enrollment_subjects(enrollment_id);
create index idx_enrollment_subjects_subject_id on enrollment_subjects(subject_id);
create index idx_enrollment_subjects_teacher_id on enrollment_subjects(teacher_id);

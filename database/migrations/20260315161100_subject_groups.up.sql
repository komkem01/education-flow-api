SET statement_timeout = 0;

--bun:split

create table subject_groups (
    id uuid primary key default gen_random_uuid(),
    school_id uuid not null references schools(id),
    code varchar(30) not null,
    name_th varchar(255) not null,
    name_en varchar(255),
    head_teacher_id uuid references member_teachers(id),
    is_active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create unique index ux_subject_groups_school_code on subject_groups(school_id, code);
create unique index ux_subject_groups_school_name_th on subject_groups(school_id, name_th);
create index idx_subject_groups_school_head_teacher on subject_groups(school_id, head_teacher_id);

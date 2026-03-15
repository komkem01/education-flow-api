SET statement_timeout = 0;

--bun:split

create table classrooms (
    id uuid primary key default gen_random_uuid(),
    school_id uuid not null references schools(id),
    academic_year_id uuid not null references academic_years(id),
    level varchar(50) not null,
    room_no varchar(20) not null,
    name varchar(120) not null,
    homeroom_teacher_id uuid references member_teachers(id),
    capacity int,
    is_active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    constraint classrooms_capacity_positive check (capacity is null or capacity > 0)
);

create unique index ux_classrooms_school_year_level_room on classrooms(school_id, academic_year_id, level, room_no);
create index idx_classrooms_school_year on classrooms(school_id, academic_year_id);

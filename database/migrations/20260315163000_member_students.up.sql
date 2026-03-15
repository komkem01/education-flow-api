SET statement_timeout = 0;

--bun:split

create table member_students (
    id uuid primary key default gen_random_uuid(),
    member_id uuid not null references members(id),
    school_id uuid not null references schools(id),
    gender_id uuid not null references genders(id),
    prefix_id uuid not null references prefixes(id),
    advisor_teacher_id uuid references member_teachers(id),
    student_code varchar(50) not null,
    first_name_th varchar(255) not null,
    last_name_th varchar(255) not null,
    first_name_en varchar(255),
    last_name_en varchar(255),
    citizen_id varchar(13),
    phone varchar(20),
    is_active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create unique index ux_member_students_school_student_code on member_students(school_id, student_code);
create unique index ux_member_students_member_id on member_students(member_id);
create unique index ux_member_students_school_citizen_id on member_students(school_id, citizen_id) where citizen_id is not null;
create index idx_member_students_member_id on member_students(member_id);
create index idx_member_students_school_id on member_students(school_id);
create index idx_member_students_advisor_teacher_id on member_students(advisor_teacher_id);
create index idx_member_students_is_active on member_students(is_active);

SET statement_timeout = 0;

--bun:split

create table school_departments (
    id uuid primary key default gen_random_uuid(),
    school_id uuid not null references schools(id),
    department_id uuid not null references departments(id),
    custom_name varchar(255),
    is_active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create unique index ux_school_departments_school_department
    on school_departments (school_id, department_id)
    where deleted_at is null;

create index idx_school_departments_school_id on school_departments (school_id);

--bun:split

insert into school_departments (school_id, department_id, is_active)
select distinct m.school_id, mm.department_id, true
from member_managements mm
join members m on m.id = mm.member_id
where mm.deleted_at is null;

--bun:split

alter table member_managements
    add column school_department_id uuid;

update member_managements mm
set school_department_id = sd.id
from members m
join school_departments sd
  on sd.school_id = m.school_id
 and sd.deleted_at is null
where m.id = mm.member_id
    and sd.department_id = mm.department_id;

alter table member_managements
    alter column school_department_id set not null,
    add constraint fk_member_managements_school_department
        foreign key (school_department_id) references school_departments(id);

create index idx_member_managements_school_department_id on member_managements (school_department_id);

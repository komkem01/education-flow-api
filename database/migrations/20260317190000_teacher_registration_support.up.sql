SET statement_timeout = 0;

--bun:split

create table if not exists teacher_addresses (
    id uuid primary key default gen_random_uuid(),
    member_teacher_id uuid not null references member_teachers(id),
    house_no varchar(50) not null,
    village varchar(100),
    road varchar(255),
    province varchar(150) not null,
    district varchar(150) not null,
    subdistrict varchar(150) not null,
    postal_code varchar(10) not null,
    is_primary boolean not null default false,
    sort_order int not null default 0,
    is_active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create index if not exists idx_teacher_addresses_teacher_id on teacher_addresses(member_teacher_id);
create index if not exists idx_teacher_addresses_sort_order on teacher_addresses(member_teacher_id, sort_order);

alter table member_teachers drop constraint if exists member_teachers_department_fkey;
alter table member_teachers
    add constraint member_teachers_department_fkey
    foreign key (department) references subject_groups(id) not valid;

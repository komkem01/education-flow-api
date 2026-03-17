SET statement_timeout = 0;

--bun:split

alter table member_teachers drop constraint if exists member_teachers_department_fkey;
alter table member_teachers
    add constraint member_teachers_department_fkey
    foreign key (department) references departments(id) not valid;

drop index if exists idx_teacher_addresses_sort_order;
drop index if exists idx_teacher_addresses_teacher_id;
drop table if exists teacher_addresses;

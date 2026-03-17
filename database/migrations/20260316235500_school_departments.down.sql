SET statement_timeout = 0;

--bun:split

alter table member_managements
    drop constraint if exists fk_member_managements_school_department;

drop index if exists idx_member_managements_school_department_id;

alter table member_managements
    drop column if exists school_department_id;

--bun:split

drop index if exists idx_school_departments_school_id;
drop index if exists ux_school_departments_school_department;
drop table if exists school_departments;

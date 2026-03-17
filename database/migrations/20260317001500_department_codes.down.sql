SET statement_timeout = 0;

--bun:split

drop index if exists ux_school_departments_school_code_active;

alter table school_departments
    drop column if exists code;

--bun:split

drop index if exists ux_departments_code_active;

alter table departments
    drop column if exists code;

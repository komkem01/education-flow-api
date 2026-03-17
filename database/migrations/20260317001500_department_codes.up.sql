SET statement_timeout = 0;

--bun:split

alter table departments
    add column if not exists code varchar(50);

update departments
set code = upper('DEP-' || substring(md5(id::text), 1, 8))
where code is null or trim(code) = '';

alter table departments
    alter column code set not null;

create unique index if not exists ux_departments_code_active
    on departments (code)
    where deleted_at is null;

--bun:split

alter table school_departments
    add column if not exists code varchar(50);

update school_departments sd
set code = upper('SD-' || substring(md5(sd.id::text), 1, 8))
from departments d
where d.id = sd.department_id
  and (sd.code is null or trim(sd.code) = '');

alter table school_departments
    alter column code set not null;

create unique index if not exists ux_school_departments_school_code_active
    on school_departments (school_id, code)
    where deleted_at is null;

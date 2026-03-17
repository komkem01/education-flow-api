SET statement_timeout = 0;

--bun:split

drop index if exists idx_member_managements_phone;

alter table member_managements
    drop column if exists phone;

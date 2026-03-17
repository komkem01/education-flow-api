SET statement_timeout = 0;

--bun:split

alter table member_managements
    add column if not exists phone varchar(20);

create index if not exists idx_member_managements_phone on member_managements (phone);

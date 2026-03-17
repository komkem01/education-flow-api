SET statement_timeout = 0;

--bun:split

alter table members
    drop constraint if exists members_email_key;

create unique index if not exists ux_members_email_active
    on members (email)
    where deleted_at is null;
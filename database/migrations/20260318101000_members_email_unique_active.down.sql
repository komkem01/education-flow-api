SET statement_timeout = 0;

--bun:split

drop index if exists ux_members_email_active;

alter table members
    add constraint members_email_key unique (email);
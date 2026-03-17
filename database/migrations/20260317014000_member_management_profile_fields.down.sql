SET statement_timeout = 0;

--bun:split

drop index if exists idx_member_managements_prefix_id;
drop index if exists idx_member_managements_gender_id;

alter table member_managements
    drop column if exists last_name,
    drop column if exists first_name,
    drop column if exists prefix_id,
    drop column if exists gender_id;

SET statement_timeout = 0;

--bun:split

alter table member_managements
    add column if not exists gender_id uuid references genders(id),
    add column if not exists prefix_id uuid references prefixes(id),
    add column if not exists first_name varchar(255),
    add column if not exists last_name varchar(255);

create index if not exists idx_member_managements_gender_id on member_managements (gender_id);
create index if not exists idx_member_managements_prefix_id on member_managements (prefix_id);

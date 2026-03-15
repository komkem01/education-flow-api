SET statement_timeout = 0;

--bun:split

create table member_guardians (
    id uuid primary key default gen_random_uuid(),
    member_id uuid references members(id),
    school_id uuid not null references schools(id),
    gender_id uuid not null references genders(id),
    prefix_id uuid not null references prefixes(id),
    first_name_th varchar(255) not null,
    last_name_th varchar(255) not null,
    first_name_en varchar(255),
    last_name_en varchar(255),
    citizen_id varchar(13),
    phone varchar(20),
    is_active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create unique index ux_member_guardians_school_citizen_id on member_guardians(school_id, citizen_id) where citizen_id is not null;
create index idx_member_guardians_member_id on member_guardians(member_id);
create index idx_member_guardians_school_id on member_guardians(school_id);
create index idx_member_guardians_is_active on member_guardians(is_active);

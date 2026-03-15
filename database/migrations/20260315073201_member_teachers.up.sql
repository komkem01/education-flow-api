SET statement_timeout = 0;

--bun:split

create table member_teachers (
    id uuid primary key default gen_random_uuid(),
    member_id uuid not null references members(id),
    gender_id uuid not null references genders(id),
    prefix_id uuid not null references prefixes(id),
    code varchar(255) not null,
    citizen_id varchar(13) not null,
    first_name_th varchar(255) not null,
    last_name_th varchar(255) not null,
    first_name_en varchar(255) not null,
    last_name_en varchar(255) not null,
    phone varchar(20) not null,
    position varchar(255) not null,
    academic_standing varchar(255) not null,
    department uuid not null references departments(id),
    start_date date not null,
    end_date date,
    is_active boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

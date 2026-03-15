SET statement_timeout = 0;

--bun:split

create table schools (
    id uuid primary key default gen_random_uuid(),
    name text not null unique,
    logo_url text,
    theme_color varchar(7),
    address text,
    description text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
)

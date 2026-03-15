SET statement_timeout = 0;

--bun:split

create table auth_sessions (
    id uuid primary key default gen_random_uuid(),
    member_id uuid not null references members(id) on delete cascade,
    token text not null unique,
    refresh_token text not null unique,
    expire_at timestamptz not null,
    refresh_expire_at timestamptz not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index idx_auth_sessions_member_id on auth_sessions(member_id);
create index idx_auth_sessions_expire_at on auth_sessions(expire_at);
create index idx_auth_sessions_refresh_expire_at on auth_sessions(refresh_expire_at);

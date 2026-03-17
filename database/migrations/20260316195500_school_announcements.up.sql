SET statement_timeout = 0;

--bun:split

create table school_announcements (
    id uuid primary key default gen_random_uuid(),
    school_id uuid not null references schools(id),
    author_member_id uuid not null references members(id),
    title text,
    content text,
    category text,
    status varchar(20) not null default 'draft',
    announced_at timestamptz,
    published_at timestamptz,
    expires_at timestamptz,
    created_by_name text,
    target_role varchar(20),
    is_pinned boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    constraint chk_school_announcements_status check (status in ('draft', 'published', 'expired')),
    constraint chk_school_announcements_target_role check (target_role is null or target_role in ('admin', 'staff', 'teacher', 'student', 'parent'))
);

create index idx_school_announcements_school_created on school_announcements (school_id, created_at desc);
create index idx_school_announcements_status on school_announcements (status);
create index idx_school_announcements_target_role on school_announcements (target_role);
create index idx_school_announcements_is_pinned on school_announcements (is_pinned);

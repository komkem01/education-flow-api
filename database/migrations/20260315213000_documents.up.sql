SET statement_timeout = 0;

--bun:split

create table documents (
    id uuid primary key default gen_random_uuid(),
    school_id uuid not null references schools(id),
    owner_member_id uuid references members(id),
    uploaded_by_member_id uuid not null references members(id),
    bucket_name text not null,
    object_key text not null,
    file_name text not null,
    content_type text not null,
    size_bytes bigint not null default 0 check (size_bytes >= 0),
    status varchar(32) not null default 'pending_upload',
    metadata jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (school_id, bucket_name, object_key)
);

create index idx_documents_school_id on documents(school_id);
create index idx_documents_owner_member_id on documents(owner_member_id);
create index idx_documents_status on documents(status);
create index idx_documents_created_at on documents(created_at);

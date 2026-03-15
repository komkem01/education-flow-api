SET statement_timeout = 0;

--bun:split

create table storages (
    id uuid primary key default gen_random_uuid(),
    school_id uuid not null references schools(id),
    provider varchar(32) not null default 's3',
    name varchar(120) not null,
    endpoint text,
    bucket_name text not null,
    is_default boolean not null default false,
    config jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz,
    unique (school_id, name),
    unique (school_id, bucket_name)
);

create unique index uq_storages_default_per_school on storages (school_id)
where is_default = true and deleted_at is null;

alter table documents add column storage_id uuid;

insert into storages (school_id, provider, name, bucket_name, is_default, created_at, updated_at)
select d.school_id,
       's3' as provider,
       d.bucket_name as name,
       d.bucket_name,
       false,
       now(),
       now()
from documents d
where d.deleted_at is null
group by d.school_id, d.bucket_name
on conflict (school_id, bucket_name) do nothing;

update documents d
set storage_id = s.id
from storages s
where d.storage_id is null
  and d.school_id = s.school_id
  and d.bucket_name = s.bucket_name;

alter table documents alter column storage_id set not null;

alter table documents
    add constraint fk_documents_storage
    foreign key (storage_id) references storages(id);

create index idx_documents_storage_id on documents(storage_id);

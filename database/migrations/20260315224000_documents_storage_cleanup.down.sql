SET statement_timeout = 0;

--bun:split

alter table documents
    add column if not exists bucket_name text;

update documents d
set bucket_name = s.bucket_name
from storages s
where d.storage_id = s.id
  and d.school_id = s.school_id
  and d.bucket_name is null;

alter table documents
    alter column bucket_name set not null;

alter table documents
    drop constraint if exists fk_documents_storage_school;

alter table documents
    add constraint fk_documents_storage
    foreign key (storage_id) references storages(id);

alter table documents
    drop constraint if exists uq_documents_school_storage_object;

alter table documents
    add constraint documents_school_id_bucket_name_object_key_key unique (school_id, bucket_name, object_key);

alter table storages
    drop constraint if exists uq_storages_id_school;

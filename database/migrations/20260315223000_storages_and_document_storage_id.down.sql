SET statement_timeout = 0;

--bun:split

drop index if exists idx_documents_storage_id;

alter table documents
    drop constraint if exists fk_documents_storage;

alter table documents
    drop column if exists storage_id;

drop index if exists uq_storages_default_per_school;
drop table if exists storages;

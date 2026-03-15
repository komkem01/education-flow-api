SET statement_timeout = 0;

--bun:split

alter table storages
    add constraint uq_storages_id_school unique (id, school_id);

alter table documents
    add constraint uq_documents_school_storage_object unique (school_id, storage_id, object_key);

alter table documents
    drop constraint if exists fk_documents_storage;

alter table documents
    add constraint fk_documents_storage_school
    foreign key (storage_id, school_id) references storages(id, school_id);

alter table documents
    drop column if exists bucket_name;

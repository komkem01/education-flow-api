SET statement_timeout = 0;

--bun:split

create unique index if not exists uq_teacher_addresses_primary_per_teacher
    on teacher_addresses(member_teacher_id)
    where is_primary = true and deleted_at is null;

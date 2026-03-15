SET statement_timeout = 0;

--bun:split

drop index if exists uk_approval_actions_idempotency;

alter table approval_actions
    drop column if exists idempotency_key;

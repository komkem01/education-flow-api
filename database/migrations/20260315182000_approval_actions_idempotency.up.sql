SET statement_timeout = 0;

--bun:split

alter table approval_actions
    add column idempotency_key varchar(120);

create unique index uk_approval_actions_idempotency
    on approval_actions(request_id, action, idempotency_key)
    where idempotency_key is not null;

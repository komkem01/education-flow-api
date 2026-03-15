SET statement_timeout = 0;

--bun:split

create table audit_logs (
    id uuid primary key default gen_random_uuid(),
    actor_id uuid,
    actor_role varchar(50),
    method varchar(10) not null,
    path text not null,
    route_path text,
    status_code int not null,
    latency_ms bigint not null,
    ip varchar(64),
    user_agent text,
    query_string text,
    request_body text,
    response_body text,
    error_message text,
    trace_id varchar(64),
    created_at timestamptz not null default now()
);

create index idx_audit_logs_actor_id on audit_logs(actor_id);
create index idx_audit_logs_method on audit_logs(method);
create index idx_audit_logs_path on audit_logs(path);
create index idx_audit_logs_status_code on audit_logs(status_code);
create index idx_audit_logs_created_at on audit_logs(created_at);

SET statement_timeout = 0;

--bun:split

create type approval_action_type as enum (
    'submit',
    'approve',
    'reject',
    'cancel',
    'comment'
);

create table approval_actions (
    id uuid primary key default gen_random_uuid(),
    request_id uuid not null references approval_requests(id) on delete cascade,
    action approval_action_type not null,
    acted_by uuid not null references members(id),
    acted_by_role approval_actor_role not null,
    comment text,
    metadata jsonb,
    created_at timestamptz not null default now()
);

create index idx_approval_actions_request_id_created_at on approval_actions(request_id, created_at desc);
create index idx_approval_actions_acted_by on approval_actions(acted_by);

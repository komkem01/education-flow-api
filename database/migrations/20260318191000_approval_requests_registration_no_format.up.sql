SET statement_timeout = 0;

create sequence if not exists approval_requests_registration_no_seq;

with numbered as (
    select id, row_number() over (order by created_at asc, id asc) as rn
    from approval_requests
)
update approval_requests ar
set registration_no = concat('APR-', lpad(numbered.rn::text, 6, '0'))
from numbered
where ar.id = numbered.id;

select setval(
    'approval_requests_registration_no_seq',
    coalesce((
        select max(substring(registration_no from '^APR-([0-9]+)$')::bigint)
        from approval_requests
    ), 0),
    true
);

alter table approval_requests
    alter column registration_no set default concat('APR-', lpad(nextval('approval_requests_registration_no_seq')::text, 6, '0'));

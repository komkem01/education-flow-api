SET statement_timeout = 0;

alter table approval_requests
    add column if not exists registration_no varchar(32);

with numbered as (
    select id, row_number() over (order by created_at asc, id asc) as rn
    from approval_requests
)
update approval_requests ar
set registration_no = concat('APR-', to_char(ar.created_at at time zone 'UTC', 'YYYYMMDD'), '-', lpad(numbered.rn::text, 8, '0'))
from numbered
where ar.id = numbered.id
  and (ar.registration_no is null or btrim(ar.registration_no) = '');

alter table approval_requests
    alter column registration_no set not null;

create unique index if not exists uq_approval_requests_registration_no on approval_requests(registration_no);

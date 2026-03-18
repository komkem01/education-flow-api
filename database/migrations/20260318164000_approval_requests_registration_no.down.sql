SET statement_timeout = 0;

drop index if exists uq_approval_requests_registration_no;

alter table approval_requests
    drop column if exists registration_no;

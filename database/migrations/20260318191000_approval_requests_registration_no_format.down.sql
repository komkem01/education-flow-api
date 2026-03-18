SET statement_timeout = 0;

alter table approval_requests
    alter column registration_no drop default;

drop sequence if exists approval_requests_registration_no_seq;

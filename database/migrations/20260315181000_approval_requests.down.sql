SET statement_timeout = 0;

--bun:split

drop table if exists approval_requests;
drop type if exists approval_request_status;
drop type if exists approval_actor_role;

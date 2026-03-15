SET statement_timeout = 0;

--bun:split

drop table if exists attendance_records;

drop type if exists attendance_source;
drop type if exists attendance_status;

SET statement_timeout = 0;

--bun:split

drop index if exists idx_teacher_experiences_teacher_sort;
drop index if exists idx_teacher_experiences_teacher_id;

--bun:split

drop table if exists teacher_experiences;

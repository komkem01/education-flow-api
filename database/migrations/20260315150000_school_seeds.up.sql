SET statement_timeout = 0;

--bun:split

insert into schools (name)
values ('โรงเรียนของนักพัฒนา')
on conflict (name) do nothing;

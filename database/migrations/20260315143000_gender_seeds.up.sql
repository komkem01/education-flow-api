SET statement_timeout = 0;

--bun:split

insert into genders (name, is_active)
values
    ('ชาย', true),
    ('หญิง', true),
    ('ไม่ระบุ', true)
on conflict (name) do nothing;

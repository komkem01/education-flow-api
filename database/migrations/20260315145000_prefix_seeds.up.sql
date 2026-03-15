SET statement_timeout = 0;

--bun:split

insert into prefixes (gender_id, name, is_active)
select g.id, v.prefix_name, true
from (
    values
        ('นาย', 'ชาย'),
        ('เด็กชาย', 'ชาย'),
        ('นาง', 'หญิง'),
        ('นางงสาว', 'หญิง'),
        ('เด็กหญิง', 'หญิง'),
        ('ไม่ระบุ', 'ไม่ระบุ')
) as v(prefix_name, gender_name)
join genders g on g.name = v.gender_name
on conflict (name) do nothing;

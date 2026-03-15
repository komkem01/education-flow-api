SET statement_timeout = 0;

--bun:split

insert into subject_groups (school_id, code, name_th, name_en, is_active)
select
    s.id,
    g.code,
    g.name_th,
    g.name_en,
    true
from schools s
cross join (
    values
        ('SG-TH', 'ภาษาไทย', 'Thai Language'),
        ('SG-MATH', 'คณิตศาสตร์', 'Mathematics'),
        ('SG-SCI', 'วิทยาศาสตร์และเทคโนโลยี', 'Science and Technology'),
        ('SG-SOC', 'สังคมศึกษา ศาสนา และวัฒนธรรม', 'Social Studies, Religion and Culture'),
        ('SG-HEALTH', 'สุขศึกษาและพลศึกษา', 'Health and Physical Education'),
        ('SG-ART', 'ศิลปะ', 'Arts'),
        ('SG-WORK', 'การงานอาชีพ', 'Occupations and Career'),
        ('SG-FLANG', 'ภาษาต่างประเทศ', 'Foreign Languages')
) as g(code, name_th, name_en)
where s.name = 'โรงเรียนของนักพัฒนา'
    and not exists (
    select 1
    from subject_groups sg
    where sg.school_id = s.id
      and (
          sg.code = g.code
          or sg.name_th = g.name_th
      )
);

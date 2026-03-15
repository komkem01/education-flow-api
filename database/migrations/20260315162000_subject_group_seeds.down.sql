SET statement_timeout = 0;

--bun:split

delete from subject_groups
where school_id in (
    select id
    from schools
    where name = 'โรงเรียนของนักพัฒนา'
)
and (code, name_th) in (
    ('SG-TH', 'ภาษาไทย'),
    ('SG-MATH', 'คณิตศาสตร์'),
    ('SG-SCI', 'วิทยาศาสตร์และเทคโนโลยี'),
    ('SG-SOC', 'สังคมศึกษา ศาสนา และวัฒนธรรม'),
    ('SG-HEALTH', 'สุขศึกษาและพลศึกษา'),
    ('SG-ART', 'ศิลปะ'),
    ('SG-WORK', 'การงานอาชีพ'),
    ('SG-FLANG', 'ภาษาต่างประเทศ')
);

SET statement_timeout = 0;

--bun:split

insert into departments (name, is_active)
values
    ('ฝ่ายบริหารวิชาการ', true),
    ('ฝ่ายบริหารงบประมาณ', true),
    ('ฝ่ายบริหารงานบุคคล', true),
    ('ฝ่ายบริหารทั่วไป', true),
    ('ฝ่ายกิจการนักเรียน', true),
    ('ฝ่ายแนะแนวและทะเบียน', true),
    ('ฝ่ายเทคโนโลยีสารสนเทศ', true),
    ('ฝ่ายอาคารสถานที่และซ่อมบำรุง', true),
    ('ฝ่ายพัสดุและจัดซื้อ', true),
    ('ฝ่ายการเงินและบัญชี', true),
    ('ฝ่ายประชาสัมพันธ์', true),
    ('ฝ่ายอนามัยโรงเรียน', true)
on conflict (name) do nothing;

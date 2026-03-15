SET statement_timeout = 0;

--bun:split

delete from prefixes
where name in ('นาย', 'เด็กชาย', 'นาง', 'นางงสาว', 'เด็กหญิง', 'ไม่ระบุ');

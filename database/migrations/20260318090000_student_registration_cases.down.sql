SET statement_timeout = 0;

--bun:split

drop table if exists student_registration_audit_logs;
drop table if exists student_registration_rules;
drop table if exists student_registration_documents;
drop table if exists student_registration_family_economic;
drop table if exists student_registration_student_guardians;
drop table if exists student_registration_guardians;
drop table if exists student_registration_previous_education;
drop table if exists student_registration_health;
drop table if exists student_registration_addresses;
drop table if exists student_registration_student_core;
drop table if exists student_registration_cases;

drop type if exists registration_income_bracket;
drop type if exists registration_document_type;
drop type if exists registration_address_type;
drop type if exists student_registration_case_status;
drop type if exists student_registration_type;

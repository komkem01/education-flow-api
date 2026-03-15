SET statement_timeout = 0;

--bun:split

create table teacher_emergency_contacts (
    id uuid primary key default gen_random_uuid(),
    member_teacher_id uuid not null references member_teachers(id),
    emergency_contact_name varchar(255) not null,
    relationship varchar(100) not null,
    phone_primary varchar(20) not null,
    phone_secondary varchar(20),
    can_decide_medical boolean not null default false,
    is_primary boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create index idx_teacher_emergency_contacts_member_teacher_id
    on teacher_emergency_contacts(member_teacher_id);

create unique index uq_teacher_emergency_contacts_primary_per_teacher
    on teacher_emergency_contacts(member_teacher_id)
    where is_primary = true and deleted_at is null;

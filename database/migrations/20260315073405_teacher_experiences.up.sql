SET statement_timeout = 0;

--bun:split

create table teacher_experiences (
	id uuid primary key default gen_random_uuid(),
	teacher_id uuid not null references member_teachers(id),
	school_name varchar(255) not null,
	position varchar(255) not null,
	department_name varchar(255),
	start_date date not null,
	end_date date,
	is_current boolean not null default false,
	responsibilities text,
	achievements text,
	sort_order int not null default 0,
	is_active boolean not null default true,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now(),
	deleted_at timestamptz,
	constraint teacher_experiences_date_check check (end_date is null or end_date >= start_date)
);

--bun:split

create index idx_teacher_experiences_teacher_id on teacher_experiences(teacher_id);
create index idx_teacher_experiences_teacher_sort on teacher_experiences(teacher_id, sort_order, start_date);

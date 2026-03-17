SET statement_timeout = 0;

--bun:split

ALTER TABLE school_announcements
    ALTER COLUMN status DROP DEFAULT,
    ALTER COLUMN status TYPE varchar(20) USING status::text,
    ALTER COLUMN status SET DEFAULT 'draft';

ALTER TABLE school_announcements
    ALTER COLUMN target_role TYPE varchar(20) USING target_role::text;

ALTER TABLE school_announcements
    ADD CONSTRAINT chk_school_announcements_status
        CHECK (status IN ('draft', 'published', 'expired')),
    ADD CONSTRAINT chk_school_announcements_target_role
        CHECK (target_role IS NULL OR target_role IN ('admin', 'staff', 'teacher', 'student', 'parent'));

DROP TYPE IF EXISTS school_announcement_target_role_enum;
DROP TYPE IF EXISTS school_announcement_status_enum;

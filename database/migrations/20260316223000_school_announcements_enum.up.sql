SET statement_timeout = 0;

--bun:split

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'school_announcement_status_enum') THEN
        CREATE TYPE school_announcement_status_enum AS ENUM ('draft', 'published', 'expired');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'school_announcement_target_role_enum') THEN
        CREATE TYPE school_announcement_target_role_enum AS ENUM ('admin', 'staff', 'teacher', 'student', 'parent');
    END IF;
END
$$;

ALTER TABLE school_announcements
    DROP CONSTRAINT IF EXISTS chk_school_announcements_status,
    DROP CONSTRAINT IF EXISTS chk_school_announcements_target_role;

ALTER TABLE school_announcements
    ALTER COLUMN status DROP DEFAULT,
    ALTER COLUMN status TYPE school_announcement_status_enum USING status::school_announcement_status_enum,
    ALTER COLUMN status SET DEFAULT 'draft'::school_announcement_status_enum;

ALTER TABLE school_announcements
    ALTER COLUMN target_role TYPE school_announcement_target_role_enum
    USING CASE
        WHEN target_role IS NULL THEN NULL
        ELSE target_role::school_announcement_target_role_enum
    END;

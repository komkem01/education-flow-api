ALTER TABLE enrollment_subjects
  DROP COLUMN IF EXISTS activity_score,
  DROP COLUMN IF EXISTS final_score,
  DROP COLUMN IF EXISTS midterm_score;

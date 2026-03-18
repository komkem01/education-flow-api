ALTER TABLE enrollment_subjects
  ADD COLUMN IF NOT EXISTS midterm_score NUMERIC(5,2),
  ADD COLUMN IF NOT EXISTS final_score NUMERIC(5,2),
  ADD COLUMN IF NOT EXISTS activity_score NUMERIC(5,2);

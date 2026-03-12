-- +goose Up 
CREATE TABLE IF NOT EXISTS applications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
  status TEXT NOT NULL DEFAULT 'applied',
  applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  interview_date TIMESTAMPTZ,
  offer_date TIMESTAMPTZ,
  notes TEXT,
  salary_offer TEXT,
  reminder_sent BOOLEAN NOT NULL DEFAULT false,
  follow_up_date TIMESTAMPTZ,
  
  -- Prevent duplicate applications (same user can't apply to same job twice)
  CONSTRAINT unique_user_job UNIQUE(user_id, job_id),
  
  -- Validate status values
  CONSTRAINT valid_status CHECK (status IN ('applied', 'interviewing', 'offer', 'rejected', 'accepted'))
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_applications_user_id ON applications(user_id);
CREATE INDEX IF NOT EXISTS idx_applications_job_id ON applications(job_id);
CREATE INDEX IF NOT EXISTS idx_applications_status ON applications(status);
CREATE INDEX IF NOT EXISTS idx_applications_interview_date ON applications(interview_date) WHERE interview_date IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_applications_follow_up_date ON applications(follow_up_date) WHERE follow_up_date IS NOT NULL;

-- Composite index for user's applications filtered by status
CREATE INDEX IF NOT EXISTS idx_applications_user_status ON applications(user_id, status);

-- +goose Down 
DROP INDEX IF EXISTS idx_applications_user_status;
DROP INDEX IF EXISTS idx_applications_follow_up_date;
DROP INDEX IF EXISTS idx_applications_interview_date;
DROP INDEX IF EXISTS idx_applications_status;
DROP INDEX IF EXISTS idx_applications_job_id;
DROP INDEX IF EXISTS idx_applications_user_id;
DROP TABLE IF EXISTS applications;

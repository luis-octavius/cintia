-- +goose Up 
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL, 
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'candidate',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index for email lookups (login queries)
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- +goose Down 
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;

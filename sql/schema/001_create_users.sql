-- +goose Up 
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY NOT NULL, 
  name TEXT NOT NULL, 
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  role TEXT DEFAULT 'candidate',
  created_at TIMESTAMP NOT NULL, 
  updated_at TIMESTAMP NOT NULL
);

-- +goose Down 
DROP TABLE IF EXISTS users;

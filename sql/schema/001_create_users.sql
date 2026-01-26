-- +goose Up 
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY NOT NULL, 
  name TEXT UNIQUE NOT NULL, 
  created_at TIMESTAMP NOT NULL, 
  updated_at TIMESTAMP NOT NULL,
  password_hash TEXT NOT NULL
);

-- +goose Down 
DROP TABLE IF EXISTS users;

-- +goose Up 
CREATE TABLE IF NOT EXISTS jobs (
  id UUID PRIMARY KEY NOT NULL, 
  name TEXT NOT NULL, 
  company TEXT NOT NULL, 
  description TEXT NOT NULL,
  link TEXT NOT NULL
);

-- +goose Down 
DROP TABLE IF EXISTS jobs; 

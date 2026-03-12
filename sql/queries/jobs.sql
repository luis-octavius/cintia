-- name: CreateJob :one
INSERT INTO jobs (
  title, 
  company, 
  location, 
  description, 
  salary_range, 
  requirements, 
  source, 
  link, 
  posted_date
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, title, company, location, description, salary_range, requirements, 
          source, link, posted_date, scraped_at, is_active, created_at, updated_at;

-- name: GetJobByID :one
SELECT id, title, company, location, description, salary_range, requirements, 
       source, link, posted_date, scraped_at, is_active, created_at, updated_at
FROM jobs
WHERE id = $1;

-- name: GetJobByLink :one
SELECT id, title, company, location, description, salary_range, requirements, 
       source, link, posted_date, scraped_at, is_active, created_at, updated_at
FROM jobs
WHERE link = $1;

-- name: UpdateJob :one
UPDATE jobs
SET
  title = COALESCE(sqlc.narg('title'), title),
  company = COALESCE(sqlc.narg('company'), company),
  location = COALESCE(sqlc.narg('location'), location),
  description = COALESCE(sqlc.narg('description'), description),
  salary_range = COALESCE(sqlc.narg('salary_range'), salary_range),
  requirements = COALESCE(sqlc.narg('requirements'), requirements),
  source = COALESCE(sqlc.narg('source'), source),
  link = COALESCE(sqlc.narg('link'), link),
  is_active = COALESCE(sqlc.narg('is_active'), is_active),
  posted_date = COALESCE(sqlc.narg('posted_date'), posted_date),
  updated_at = NOW()
WHERE id = $1
RETURNING id, title, company, location, description, salary_range, requirements, 
          source, link, posted_date, scraped_at, is_active, created_at, updated_at;

-- name: DeleteJob :exec
DELETE FROM jobs WHERE id = $1;

-- name: ListJobs :many
SELECT id, title, company, location, description, salary_range, requirements, 
       source, link, posted_date, scraped_at, is_active, created_at, updated_at
FROM jobs
WHERE 
  (sqlc.narg('title')::TEXT IS NULL OR title ILIKE '%' || sqlc.narg('title')::TEXT || '%')
  AND (sqlc.narg('company')::TEXT IS NULL OR company ILIKE '%' || sqlc.narg('company')::TEXT || '%')
  AND (sqlc.narg('location')::TEXT IS NULL OR location ILIKE '%' || sqlc.narg('location')::TEXT || '%')
  AND (sqlc.narg('source')::TEXT IS NULL OR source = sqlc.narg('source')::TEXT)
  AND (sqlc.narg('is_active')::BOOLEAN IS NULL OR is_active = sqlc.narg('is_active')::BOOLEAN)
ORDER BY posted_date DESC
LIMIT $1 OFFSET $2;

-- name: CountJobs :one
SELECT COUNT(*)
FROM jobs
WHERE 
  (sqlc.narg('title')::TEXT IS NULL OR title ILIKE '%' || sqlc.narg('title')::TEXT || '%')
  AND (sqlc.narg('company')::TEXT IS NULL OR company ILIKE '%' || sqlc.narg('company')::TEXT || '%')
  AND (sqlc.narg('location')::TEXT IS NULL OR location ILIKE '%' || sqlc.narg('location')::TEXT || '%')
  AND (sqlc.narg('source')::TEXT IS NULL OR source = sqlc.narg('source')::TEXT)
  AND (sqlc.narg('is_active')::BOOLEAN IS NULL OR is_active = sqlc.narg('is_active')::BOOLEAN);

-- name: CreateApplication :one
INSERT INTO applications (
  user_id,
  job_id,
  notes
)
VALUES ($1, $2, $3)
RETURNING id, user_id, job_id, status, applied_at, updated_at, 
          interview_date, offer_date, notes, salary_offer, reminder_sent, follow_up_date;

-- name: GetApplicationByID :one
SELECT id, user_id, job_id, status, applied_at, updated_at, 
       interview_date, offer_date, notes, salary_offer, reminder_sent, follow_up_date
FROM applications
WHERE id = $1;

-- name: GetUserApplications :many
SELECT id, user_id, job_id, status, applied_at, updated_at, 
       interview_date, offer_date, notes, salary_offer, reminder_sent, follow_up_date
FROM applications
WHERE user_id = $1
ORDER BY applied_at DESC;

-- name: GetJobApplications :many
SELECT id, user_id, job_id, status, applied_at, updated_at, 
       interview_date, offer_date, notes, salary_offer, reminder_sent, follow_up_date
FROM applications
WHERE job_id = $1
ORDER BY applied_at DESC;

-- name: GetApplicationByUserAndJob :one
SELECT id, user_id, job_id, status, applied_at, updated_at, 
       interview_date, offer_date, notes, salary_offer, reminder_sent, follow_up_date
FROM applications
WHERE user_id = $1 AND job_id = $2;

-- name: UpdateApplication :one
UPDATE applications
SET
  interview_date = COALESCE(sqlc.narg('interview_date'), interview_date),
  offer_date = COALESCE(sqlc.narg('offer_date'), offer_date),
  notes = COALESCE(sqlc.narg('notes'), notes),
  salary_offer = COALESCE(sqlc.narg('salary_offer'), salary_offer),
  reminder_sent = COALESCE(sqlc.narg('reminder_sent'), reminder_sent),
  follow_up_date = COALESCE(sqlc.narg('follow_up_date'), follow_up_date),
  updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, job_id, status, applied_at, updated_at, 
          interview_date, offer_date, notes, salary_offer, reminder_sent, follow_up_date;

-- name: UpdateApplicationStatus :one
UPDATE applications
SET
  status = $2,
  updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, job_id, status, applied_at, updated_at, 
          interview_date, offer_date, notes, salary_offer, reminder_sent, follow_up_date;

-- name: DeleteApplication :exec
DELETE FROM applications WHERE id = $1;

-- name: GetApplicationsByStatus :many
SELECT id, user_id, job_id, status, applied_at, updated_at, 
       interview_date, offer_date, notes, salary_offer, reminder_sent, follow_up_date
FROM applications
WHERE user_id = $1 AND status = $2
ORDER BY applied_at DESC;

-- name: GetUpcomingInterviews :many
SELECT id, user_id, job_id, status, applied_at, updated_at, 
       interview_date, offer_date, notes, salary_offer, reminder_sent, follow_up_date
FROM applications
WHERE user_id = $1 
  AND interview_date IS NOT NULL 
  AND interview_date >= NOW()
  AND status = 'interviewing'
ORDER BY interview_date ASC;

-- name: GetPendingReminders :many
SELECT id, user_id, job_id, status, applied_at, updated_at, 
       interview_date, offer_date, notes, salary_offer, reminder_sent, follow_up_date
FROM applications
WHERE reminder_sent = false
  AND interview_date IS NOT NULL
  AND interview_date <= NOW() + INTERVAL '24 hours'
  AND interview_date >= NOW();

-- name: MarkReminderSent :exec
UPDATE applications
SET reminder_sent = true, updated_at = NOW()
WHERE id = $1;

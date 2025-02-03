-- name: GetApplicationByFileID :one
SELECT a.id
FROM applications a
JOIN candidates c ON a.candidate_id = c.candidate_id
JOIN files f ON c.file_id = f.id
WHERE f.id = $1;

-- name: GetApplication :one
select *
from applications
where id = $1
limit 1
;

-- name: ListApplications :many
select *
from applications
;

-- name: CreateApplication :one
INSERT INTO applications (
    id, candidate_id, job_role_id, stage, score, applied_at, last_updated
) VALUES (
    uuid_generate_v4(),  $1, $2, $3, $4, now(), now()
)
RETURNING *;

-- name: UpdateApplication :exec
UPDATE applications
SET candidate_id = $2,
    job_role_id = $3,
    stage = $4,
    score = $5,
    last_updated = now()
WHERE id = $1;

-- name: DeleteApplication :exec
delete from applications
where id = $1
;

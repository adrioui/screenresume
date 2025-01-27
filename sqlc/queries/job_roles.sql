-- name: GetJobRole :one
select *
from job_roles
where id = $1
limit 1
;

-- name: ListJobRoles :many
select *
from job_roles
order by created_at desc
;

-- name: CreateJobRole :one
INSERT INTO job_roles (
    id, title, department_id, level, salary_range, location, is_active, created_at, updated_at
) VALUES (
    uuid_generate_v4(), $1, $2, $3, $4, $5, $6, now(), now()
)
RETURNING *;

-- name: UpdateJobRole :exec
UPDATE job_roles
SET title = $2,
    department_id = $3,
    level = $4,
    salary_range = $5,
    location = $6,
    is_active = $7,
    updated_at = now()
WHERE id = $1;

-- name: DeleteJobRole :exec
delete from job_roles
where id = $1
;

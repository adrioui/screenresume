-- name: GetJobRoleRequirement :one
select *
from job_role_requirements
where job_role_id = $1 and skill_id = $2
limit 1
;

-- name: ListJobRoleRequirements :many
select *
from job_role_requirements
order by job_role_id, skill_id
;

-- name: CreateJobRoleRequirement :one
INSERT INTO job_role_requirements (
  job_role_id, skill_id, required, min_experience_years, importance
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateJobRoleRequirement :exec
UPDATE job_role_requirements
SET required = $3,
    min_experience_years = $4,
    importance = $5
WHERE job_role_id = $1 AND skill_id = $2;

-- name: ToggleJobRoleRequirement :exec
UPDATE job_role_requirements
SET required = NOT required
WHERE job_role_id = $1 AND skill_id = $2;

-- name: DeleteJobRoleRequirement :exec
delete from job_role_requirements
where job_role_id = $1 and skill_id = $2
;

-- name: GetScreeningCriteria :one
select *
from screening_criteria
where id = $1
limit 1
;

-- name: ListScreeningCriteria :many
select *
from screening_criteria
;

-- name: CreateScreeningCriteria :one
insert into screening_criteria (
  id, screening_result_id, criteria_text, decision, reasoning, matched_skills, missing_skills
) values (
  uuid_generate_v4(), $1, $2, $3, $4, $5, $6
)
RETURNING *;

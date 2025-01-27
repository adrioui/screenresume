-- name: GetCandidateSkills :one
select *
from candidate_skills
where candidate_id = $1 and skill_id = $2
limit 1
;

-- name: ListCandidateSkills :many
select *
from candidate_skills
order by candidate_id, skill_id
;

-- name: CreateCandidateSkills :one
INSERT INTO candidate_skills (
    candidate_id, skill_id, years_experience, last_used
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateCandidateSkills :exec
UPDATE candidate_skills
SET 
    years_experience = $3, 
    last_used = $4
WHERE candidate_id = $1 AND skill_id = $2;

-- name: DeleteCandidateSkills :exec
delete from candidate_skills
where candidate_id = $1 and skill_id = $2
;

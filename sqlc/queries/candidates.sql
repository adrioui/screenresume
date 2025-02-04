-- name: GetCandidate :one
select *
from candidates
where id = $1
limit 1
;

-- name: ListCandidates :many
select *
from candidates
order by created_at desc
;

-- name: CreateCandidate :one
INSERT INTO candidates (
    id, full_name, email, phone, file_id, status
) VALUES (
    uuid_generate_v4(), $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateCandidate :exec
UPDATE candidates
SET full_name = $2,
    email = $3,
    phone = $4,
    file_id = $5,
    status = $6
WHERE id = $1;

-- name: DeleteCandidate :exec
delete from candidates
where id = $1
;

-- name: CandidateAndJobRoles :many
select sqlc.embed(c), sqlc.embed(jr), a.applied_at
from candidates c
join applications a on c.id = a.candidate_id
join job_roles jr on a.job_role_id = jr.id
where jr.is_active = true
and (c.full_name ilike '%' || sqlc.arg(name_search) || '%' or sqlc.arg(name_search) is null)
order by case when c.full_name = sqlc.arg(name_search) then 0 else 1 end, c.full_name asc
limit sqlc.arg(limitQuery) offset sqlc.arg(pageQuery)
;

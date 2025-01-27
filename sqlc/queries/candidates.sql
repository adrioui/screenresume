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

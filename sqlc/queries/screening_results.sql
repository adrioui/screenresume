-- name: GetScreeningResults :one
select *
from screening_results
where id = $1
limit 1
;

-- name: ListScreeningResults :many
select *
from screening_results
order by id
;

-- name: CreateScreeningResults :one
INSERT INTO screening_results (
    id, application_id, model_version, raw_response, processed_at
) VALUES (
    uuid_generate_v4(), $1, $2, $3, now()
)
RETURNING *;
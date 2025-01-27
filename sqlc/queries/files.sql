-- name: GetFile :one
select *
from files
where id = $1
limit 1
;

-- name: GetFilesByPath :many
select *
from files
where path = $1
order by created_at
;

-- name: ListFiles :many
select *
from files
order by created_at
;

-- name: CreateFile :one
insert into files (
  id, path, file_type, checksum
) values (
  uuid_generate_v4(), $1, $2, $3
)
RETURNING *;

-- name: UpdateFile :exec
update files
  set path = $2,
      file_type = $3,
      checksum = $4
where id = $1;

-- name: DeleteFile :exec
delete from files
where id = $1
;

-- name: GetDepartment :one
select *
from departments
where id = $1
limit 1
;

-- name: ListDepartments :many
select *
from departments
order by created_at
;

-- name: CreateDepartment :one
insert into departments (
  id, name
) values (
  uuid_generate_v4(), $1
)
RETURNING *;

-- name: UpdateDepartment :exec
update departments
  set  name = $2
where id = $1;

-- name: DeleteDepartment :exec
delete from departments
where id = $1
;

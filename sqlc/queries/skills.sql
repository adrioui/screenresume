-- name: GetSkill :one
select *
from skills
where id = $1
limit 1
;

-- name: ListSkills :many
select *
from skills
order by created_at
;

-- name: CreateSkill :one
insert into skills (
  id, name, category
) values (
  uuid_generate_v4(), $1, $2
)
RETURNING *;

-- name: UpdateSkill :exec
update skills
  set name = $2,
      category = $3
where id = $1;

-- name: DeleteSkill :exec
delete from skills
where id = $1
;

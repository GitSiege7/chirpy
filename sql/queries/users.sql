-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, hashed_password)
values (
    gen_random_uuid(),
    now(),
    now(),
    $1, 
    $2
)
returning *;

-- name: DeleteUsers :exec
delete from users;

-- name: GetUserByEmail :one
select * from users
where users.email = $1;

-- name: UpdateCredentials :one
update users
set email = $2, hashed_password = $3
where users.id = $1
returning *;
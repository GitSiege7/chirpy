-- name: CreateChirp :one
insert into chirps (id, created_at, updated_at, body, user_id)
values (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)
returning *;

-- name: GetChirps :many
select * from chirps
order by chirps.created_at asc;

-- name: GetChirpByID :one
select * from chirps
where chirps.id = $1;

-- name: DeleteChirp :exec
delete from chirps
where chirps.id = $1;

-- name: GetChirpsByUser :many
select * from chirps
where chirps.user_id = $1
order by chirps.created_at asc;
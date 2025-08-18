-- name: CreateRefreshToken :one
insert into refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
values (
    $1,
    Now(),
    Now(),
    $2,
    Now() + interval '60 days',
    null
)
returning *;

-- name: GetRefreshToken :one
select * from refresh_tokens
where refresh_tokens.token = $1;

-- name: GetUserFromRefreshToken :one
select refresh_tokens.user_id from refresh_tokens
where refresh_tokens.token = $1;

-- name: SetRevoked :exec
update refresh_tokens
set revoked_at = Now(), updated_at = Now()
where refresh_tokens.token = $1;
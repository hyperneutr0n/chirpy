-- name: CreateUser :one
INSERT INTO users (email, password)
VALUES (
    $1,
    $2
)
RETURNING id, email, created_at, updated_at;

-- name: ResetUser :exec
DELETE FROM users;

-- name: LoginUser :one
SELECT email, password
FROM users
WHERE email=$1;

-- name: GetUserByEmail :one
SELECT id, email, created_at, updated_at
FROM users
WHERE email=$1;

-- name: GetUserFromRefreshToken :one
SELECT * 
FROM users
WHERE id=(
    SELECT user_id 
    FROM refresh_tokens 
    WHERE token=$1
    );

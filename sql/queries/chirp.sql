-- name: CreateChirp :one
INSERT INTO chirps(body, user_id)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
WHERE 
    (user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL)
ORDER BY 
    CASE WHEN @sort_dir::text = 'desc' THEN created_at END DESC,
    CASE WHEN @sort_dir::text != 'desc' THEN created_at END ASC;

-- name: FindChirp :one
SELECT * FROM chirps WHERE id=$1;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id=$1;
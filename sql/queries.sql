-- name: CreateArt :one
INSERT INTO art (name, bio)
VALUES ($1, $2)
RETURNING *;

-- name: GetArt :one
SELECT *
FROM art
WHERE id = $1
LIMIT 1;

-- name: UpdateArt :one
UPDATE art
SET name = $2,
    bio  = $3
WHERE id = $1
RETURNING *;

-- name: PartialUpdateArt :one
UPDATE art
SET name = CASE WHEN @update_name::boolean THEN @name::VARCHAR(32) ELSE name END,
    bio  = CASE WHEN @update_bio::boolean THEN @bio::TEXT ELSE bio END
WHERE id = @id
RETURNING *;

-- name: DeleteArt :exec
DELETE
FROM art
WHERE id = $1;

-- name: ListArt :many
SELECT *
FROM art
ORDER BY name;

-- name: TruncateArt :exec
TRUNCATE art;
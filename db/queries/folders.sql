-- GET queries
-- name: FindFolderById :one
SELECT *
FROM folders 
WHERE id = $1;
-- name: FindFolderByUserIdMostRecent :one
SELECT *
FROM folders 
WHERE user_id = $1
ORDER BY updated_datetime DESC
LIMIT 1;
-- name: FindFoldersByUserId :many
SELECT *
FROM folders 
WHERE user_id = $1;
-- name: FindFoldersByParentId :many
SELECT *
FROM folders 
WHERE parent_id = $1;
-- name: FindFoldersAll :many
SELECT *
FROM folders;

-- POST queries
-- name: CreateFolder :one
INSERT INTO folders (
    user_id, parent_id, name, description
) VALUES ( $1, $2, $3, $4 )
RETURNING *;

-- PUT queries
-- name: UpdateFolder :exec
UPDATE folders 
SET name = $1, description = $2, updated_datetime = now()
WHERE id = $3;
-- name: MoveFolder :exec
UPDATE folders 
SET parent_id = $2, updated_datetime = now()
WHERE id = $1;

-- PATCH queries
-- name: UpdateFolderDateTime :exec
UPDATE folders 
SET updated_datetime = now()
WHERE id = $1;


-- Delete queries
-- name: DeleteFolder :exec
DELETE FROM folders WHERE id = $1;



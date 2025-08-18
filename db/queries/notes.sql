-- GET queries
-- name: FindNoteById :one
SELECT *
FROM notes 
WHERE id = $1;
-- name: FindNotesByFolderId :many
SELECT *
FROM notes 
WHERE folder_id = $1;
-- name: FindNotesByUserId :many
SELECT *
FROM notes 
WHERE user_id = $1;
-- name: FindNotesAll :many
SELECT *
FROM notes;
-- name: FindNoteByUserIdMostRecent :one
SELECT *
FROM notes 
WHERE user_id = $1
ORDER BY updated_datetime DESC
LIMIT 1;
-- name: FindNotesByUserIdMostRecent :many
SELECT *
FROM notes 
WHERE user_id = $1
ORDER BY updated_datetime DESC;
-- name: FindNotesByUserIDDateTimeRange :many
SELECT *
FROM notes 
WHERE user_id = $1
AND created_datetime >= $2
AND created_datetime <= $3;
-- name: FindNotesByUserIDDateTimeRangeUpdated :many
SELECT *
FROM notes 
WHERE user_id = $1
AND updated_datetime >= $2
AND updated_datetime <= $3;

-- POST queries
-- name: CreateNote :one
INSERT INTO notes (
    user_id, folder_id, name, description, content
) VALUES ( $1, $2, $3, $4, $5 )
RETURNING *;

-- PATCH queries
-- name: UpdateNote :exec
UPDATE notes 
SET name = $1, description = $2, content = $3, updated_datetime = now()
WHERE id = $4;
-- name: MoveNote :exec
-- also update updated_datetime for previous folder
UPDATE notes 
SET folder_id = $2, updated_datetime = now()
WHERE id = $1;
UPDATE notes 
SET updated_datetime = now()
WHERE folder_id = (
    SELECT folder_id FROM notes WHERE id = $1
);
-- DELETE queries
-- name: DeleteNote :exec
DELETE FROM notes WHERE id = $1;

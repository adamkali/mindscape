-- name: FindBookmarkById :one
SELECT *
FROM bookmarks 
WHERE id = $1;
-- name: FindBookmarksByFolderId :many
SELECT *
FROM bookmarks 
WHERE folder_id = $1;
-- name: FindBookmarksByUserId :many
SELECT *
FROM bookmarks 
WHERE user_id = $1;
-- name: FindBookmarksAll :many
SELECT *
FROM bookmarks;
-- name: FindBookmarkByUserIdMostRecent :one
SELECT *
FROM bookmarks 
WHERE user_id = $1
ORDER BY updated_datetime DESC
LIMIT 1;
-- name: FindBookmarksByUserIdMostRecent :many
SELECT *
FROM bookmarks 
WHERE user_id = $1
ORDER BY updated_datetime DESC;
-- name: FindBookmarksByUserIDDateTimeRange :many
SELECT *
FROM bookmarks 
WHERE user_id = $1
AND updated_datetime BETWEEN $2 AND $3;

-- name: CreateBookmark :one
INSERT INTO bookmarks (
    user_id, folder_id, name, link
) VALUES ( $1, $2, $3, $4 )
RETURNING *;

-- name: UpdateBookmark :exec
UPDATE bookmarks 
SET folder_id = $2, name = $3, link = $4, updated_datetime = now()
WHERE id = $1;
-- name: MoveBookmark :exec
UPDATE bookmarks 
SET folder_id = $2, updated_datetime = now()
WHERE id = $1;
UPDATE bookmarks 
SET updated_datetime = now()
WHERE folder_id = (SELECT folder_id FROM bookmarks WHERE id = $1);

-- name: DeleteBookmark :exec
DELETE FROM bookmarks WHERE id = $1;

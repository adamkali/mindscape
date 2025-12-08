-- name: FindUserWidgetsByUserID :many
SELECT * FROM user_widgets WHERE user_id = $1;
-- name: FindUserWidgetByID :one
SELECT * FROM user_widgets WHERE id = $1;

-- name: CreateUserWidget :one
INSERT INTO user_widgets (
	user_id,
	schema_id,
	config,
	position_x,
	position_y,
	width,
	height,
	z_index,
	is_visible
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: UpdateUserWidget :one
UPDATE user_widgets SET
	config = $2,
	position_x = $3,
	position_y = $4,
	width = $5,
	height = $6,
	z_index = $7,
	is_visible = $8
WHERE id = $1 RETURNING *;

-- name: DeleteUserWidget :exec
DELETE FROM user_widgets WHERE id = $1;

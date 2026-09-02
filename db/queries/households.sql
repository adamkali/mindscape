-- GET queries
-- name: GetHouseholdByID :one
SELECT *
FROM households
WHERE id = $1;

-- name: GetHouseholdsByUserID :many
SELECT h.*
FROM households h
JOIN household_members m ON m.household_id = h.id
WHERE m.user_id = $1
ORDER BY h.created_at DESC;

-- name: GetMembersByHouseholdID :many
SELECT m.id, m.household_id, m.user_id, m.role, m.joined_at, u.username
FROM household_members m
JOIN users u ON u.id = m.user_id
WHERE m.household_id = $1
ORDER BY m.joined_at ASC;

-- name: GetMember :one
SELECT *
FROM household_members
WHERE household_id = $1 AND user_id = $2;

-- name: GetInviteByCode :one
SELECT *
FROM household_invites
WHERE code = $1;

-- name: GetInvitesByUserID :many
SELECT i.id, i.household_id, i.invited_user_id, i.code, i.expires_at, i.created_by, i.created_at, h.name AS household_name
FROM household_invites i
JOIN households h ON h.id = i.household_id
WHERE i.invited_user_id = $1 AND i.accepted_at IS NULL AND i.expires_at > now()
ORDER BY i.created_at DESC;

-- POST queries
-- name: CreateHousehold :one
INSERT INTO households (
    name, owner_id
) VALUES ( $1, $2 )
RETURNING *;

-- name: AddMember :one
INSERT INTO household_members (
    household_id, user_id, role
) VALUES ( $1, $2, $3 )
RETURNING *;

-- name: CreateInvite :one
INSERT INTO household_invites (
    household_id, invited_user_id, code, expires_at, created_by
) VALUES ( $1, $2, $3, $4, $5 )
RETURNING *;

-- PUT queries
-- name: MarkInviteAccepted :exec
UPDATE household_invites
SET accepted_at = now()
WHERE id = $1;

-- DELETE queries
-- name: RemoveMember :exec
DELETE FROM household_members
WHERE household_id = $1 AND user_id = $2;

-- name: DeleteInvite :exec
DELETE FROM household_invites WHERE id = $1;

-- name: DeleteHousehold :exec
DELETE FROM households WHERE id = $1;

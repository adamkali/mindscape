-- name: GetTaskById :one
SELECT *
FROM tasks 
WHERE id = $1;

-- name: GetTasksByUserId :many
SELECT *
FROM tasks 
WHERE user_id = $1;

-- name: GetTasks :many
SELECT *
FROM tasks;

-- name: GetTasksByScheduledTaskType :many
SELECT *
FROM tasks
WHERE task_type_id IN (
	SELECT id FROM task_types WHERE show_in_scheduled = true
) AND user_id = $1
ORDER BY updated_at DESC, due_at DESC;

-- name: GetTasksByCompletedTaskType :many
SELECT *
FROM tasks
WHERE task_type_id IN (
	SELECT id FROM task_types WHERE show_in_completed = true
) AND user_id = $1
ORDER BY updated_at DESC, completed_at DESC;

-- name: GetTasksByAvailableTaskType :many
SELECT *
FROM tasks
WHERE task_type_id IN (
	SELECT id FROM task_types WHERE show_in_available = true
) AND user_id = $1
ORDER BY created_at DESC, due_at DESC;

-- name: GetTasksByCancelledTaskType :many
SELECT *
FROM tasks
WHERE task_type_id IN (
	SELECT id FROM task_types WHERE show_in_cancelled = true
) AND user_id = $1
ORDER BY created_at DESC, completed_at DESC;


-- name: GetTasksByTaskType :many
SELECT *
FROM tasks 
WHERE task_type_id = $1 AND user_id = $2
ORDER BY updated_at DESC;


-- InsertNewTask
-- name: InsertNewTask :one
INSERT INTO tasks (user_id, task_type_id, name, description, due_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- UpdateAsAmbiguous
-- name: UpdateAsAmbiguous :one
UPDATE tasks
SET task_type_id = (SELECT id FROM task_types WHERE name = 'AmbiguousTaskStatus'),
    completed_at = NULL,
    due_at = NULL,
    updated_at = now()
WHERE tasks.id = $1
RETURNING *;

-- UpdateAsCancelled
-- name: UpdateAsCancelled :one
UPDATE tasks
SET task_type_id = (SELECT id FROM task_types WHERE name = 'CancelledTaskStatus'),
    completed_at = now(),
    due_at = NULL,
    updated_at = now()
WHERE tasks.id = $1
RETURNING *;


-- UpdateAsDone
-- name: UpdateAsDone :one
UPDATE tasks
SET task_type_id = (SELECT id FROM task_types WHERE name = 'DoneTaskStatus'),
    completed_at = now(),
    due_at = NULL,
    updated_at = now()
WHERE tasks.id = $1
RETURNING *;


-- UpdateAsHold
-- name: UpdateAsHold :one
UPDATE tasks
SET task_type_id = (SELECT id FROM task_types WHERE name = 'HoldTaskStatus'),
    completed_at = NULL,
    due_at = NULL,
    updated_at = now()
WHERE tasks.id = $1
RETURNING *;


-- UpdateAsPending
-- name: UpdateAsPending :one
UPDATE tasks
SET task_type_id = (SELECT id FROM task_types WHERE name = 'PendingTaskStatus'),
    completed_at = NULL,
    due_at = $2,
    updated_at = now()
WHERE tasks.id = $1
RETURNING *;


-- UpdateAsRecurring
-- name: UpdateAsRecurring :one
UPDATE tasks
SET task_type_id = (SELECT id FROM task_types WHERE name = 'RecurringTaskStatus'),
    completed_at = NULL,
    due_at = now() + interval '1 days',
    updated_at = now()
WHERE tasks.id = $1
RETURNING *;

-- UdateAsUndone
-- name: UpdateAsUndone :one
UPDATE tasks
SET task_type_id = (SELECT id FROM task_types WHERE name = 'UndoneTaskStatus'),
    completed_at = NULL,
    due_at = $2,
    updated_at = now()
WHERE tasks.id = $1
RETURNING *;

--UpdateAsUrgent
-- name: UpdateAsUrgent :one
UPDATE tasks
SET task_type_id = (SELECT id FROM task_types WHERE name = 'UrgentTaskStatus'),
    completed_at = NULL,
    due_at = now() + interval '4 hours',
    updated_at = now()
WHERE tasks.id = $1
RETURNING *;


--DeleteTask
-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;


--UpdateTaskContent
-- name: UpdateTaskContent :one
UPDATE tasks
SET updated_at = now(),
    name = $2,
    description = $3,
    due_at = $4
WHERE id = $1
RETURNING *;


-- GET queries
-- name: GetTaskTypeById :one
SELECT *
FROM task_types 
WHERE id = $1;

-- name: GetTaskTypes :many
SELECT *
FROM task_types;

-- name: GetTaskTypeByName :one
SELECT *
FROM task_types 
WHERE name = $1;


-- name: GetAmbiguousTaskType :one
SELECT *
FROM task_types 
WHERE name = 'AmbiguousTaskStatus';

-- name: GetCancelledTaskType :one
SELECT *
FROM task_types 
WHERE name = 'CancelledTaskStatus';

-- name: GetDoneTaskType :one
SELECT *
FROM task_types 
WHERE name = 'DoneTaskStatus';


-- name: GetHoldTaskType :one
SELECT *
FROM task_types 
WHERE name = 'HoldTaskStatus';

-- name: GetPendingTaskType :one
SELECT *
FROM task_types 
WHERE name = 'PendingTaskStatus';


-- name: GetRecurringTaskType :one
SELECT *
FROM task_types 
WHERE name = 'RecurringTaskStatus';


-- name: GetUndoneTaskType :one
SELECT *
FROM task_types 
WHERE name = 'UndoneTaskStatus';

-- name: GetUrgentTaskType :one
SELECT *
FROM task_types 
WHERE name = 'UrgentTaskStatus';




package requests

import "github.com/google/uuid"

type CreateTaskRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TaskTypeID  uuid.UUID `json:"task_type_id"`
} // @name CreateTaskRequest

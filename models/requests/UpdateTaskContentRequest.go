package requests

import "github.com/google/uuid"

type UpdateTaskContentRequest struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
} // @name UpdateTaskContentRequest

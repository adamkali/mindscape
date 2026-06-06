package requests

import (
	"time"

	"github.com/google/uuid"
)

type UpdateTaskContentRequest struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description" validate:"required"`
	DueAt       *time.Time `json:"due_at"`
} // @name UpdateTaskContentRequest

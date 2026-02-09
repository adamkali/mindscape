package services

import (
	"time"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
)

type TaskDTO struct {
	ID          uuid.UUID           `json:"id"`
	UserID      uuid.UUID           `json:"user_id"`
	TaskTypeID  uuid.UUID           `json:"task_type_id"`
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	CreatedAt   time.Time           `json:"created_at"`
	DueAt       time.Time           `json:"due_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	CompletedAt time.Time           `json:"completed_at"`
	TaskType    repository.TaskType `json:"task_type"`
}

type ITaskService interface {
	GetAll(userId uuid.UUID) ([]TaskDTO, error)
	GetById(userId uuid.UUID, taskId uuid.UUID) (TaskDTO, error)
	GetByTaskByUserID(userId uuid.UUID) ([]TaskDTO, error)
	GetTasksByScheduledTaskType(userId uuid.UUID) ([]TaskDTO, error)
	GetTasksByCompletedTaskType(userId uuid.UUID) ([]TaskDTO, error)
	GetTasksByAvailableTaskType(userId uuid.UUID) ([]TaskDTO, error)
	GetTasksByCancelledTaskType(userId uuid.UUID) ([]TaskDTO, error)
	GetTasksByTaskType(userId uuid.UUID, taskType string) ([]TaskDTO, error)
	Create(repository.InsertNewTaskParams) (TaskDTO, error)
	Delete( taskId uuid.UUID) error
	UpdateAsAmbiguous(uuid.UUID) (TaskDTO, error)
	UpdateAsCancelled(uuid.UUID) (TaskDTO, error)
	UpdateAsDone(uuid.UUID) (TaskDTO, error)
	UpdateAsHold(uuid.UUID) (TaskDTO, error)
	UpdateAsPending(repository.UpdateAsPendingParams) (TaskDTO, error)
	UpdateAsRecurring(uuid.UUID) (TaskDTO, error)
	UpdateAsUndone(repository.UpdateAsUndoneParams) (TaskDTO, error)
	UpdateAsUrgent(userId uuid.UUID) (TaskDTO, error)
	UpdateTaskContent(repository.UpdateTaskContentParams) (TaskDTO, error)
	DeleteTaskContent(uuid.UUID) (TaskDTO, error)
}

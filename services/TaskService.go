package services

import (
	"context"
	"time"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskService struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func CreateTaskService(
	ctx context.Context,
	pool *pgxpool.Pool,
) ITaskService {
	return &TaskService{
		ctx:  ctx,
		pool: pool,
	}
}

func (s *TaskService) getTaskTypeById(taskId uuid.UUID, repo *repository.Queries) (repository.TaskType, error) {
	task_status, err := repo.GetTaskTypeById(s.ctx, taskId)
	if err != nil {
		return repository.TaskType{}, err
	}
	return task_status, nil
}

func (s *TaskService) getTaskTypesForSlice(tasks []repository.Task, repo *repository.Queries) ([]TaskDTO, error) {
	taskTypes := make(chan repository.TaskType)
	errChans := make(chan error)
	taskObjects := make([]TaskDTO, len(tasks))
	for _, task := range tasks {
		go func() {
			taskType, err := s.getTaskTypeById(task.TaskTypeID, repo)
			if err != nil {
				errChans <- err
			}
			taskTypes <- taskType
		}()
	}
	for i := range tasks {
		select {
		case taskType := <-taskTypes:
			taskObjects[i] = TaskDTO{
				ID:          tasks[i].ID,
				UserID:      tasks[i].UserID,
				TaskTypeID:  tasks[i].TaskTypeID,
				Name:        tasks[i].Name,
				Description: tasks[i].Description,
				CreatedAt:   tasks[i].CreatedAt.Time,
				DueAt:       tasks[i].DueAt.Time,
				UpdatedAt:   tasks[i].UpdatedAt.Time,
				CompletedAt: tasks[i].CompletedAt.Time,
				TaskType:    taskType,
			}
		case err := <-errChans:
			return []TaskDTO{}, err
		}
	}

	return taskObjects, nil
}

func (s *TaskService) GetAll(userId uuid.UUID) ([]TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return []TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	tasks, err := repo.GetTasksByUserId(s.ctx, userId)
	if err != nil {
		return []TaskDTO{}, err
	}
	taskObjects, err := s.getTaskTypesForSlice(tasks, repo)
	tx.Commit(s.ctx)
	return taskObjects, nil
}

func (s *TaskService) GetById(userId uuid.UUID, taskId uuid.UUID) (TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	task, err := repo.GetTaskById(s.ctx, taskId)
	if err != nil {
		return TaskDTO{}, err
	}

	taskType, err := s.getTaskTypeById(task.TaskTypeID, repo)
	if err != nil {
		return TaskDTO{}, err
	}
	taskModel := TaskDTO{
		ID:          task.ID,
		UserID:      task.UserID,
		TaskTypeID:  task.TaskTypeID,
		Name:        task.Name,
		Description: task.Description,
		CreatedAt:   task.CreatedAt.Time,
		DueAt:       task.DueAt.Time,
		UpdatedAt:   task.UpdatedAt.Time,
		CompletedAt: task.CompletedAt.Time,
		TaskType:    taskType,
	}
	tx.Commit(s.ctx)
	return taskModel, nil
}
func (s *TaskService) GetByTaskByUserID(userId uuid.UUID) ([]TaskDTO, error)
func (s *TaskService) GetTasksByScheduledTaskType(userId uuid.UUID) ([]TaskDTO, error) 
func (s *TaskService) GetTasksByCompletedTaskType(userId uuid.UUID) ([]TaskDTO, error)
func (s *TaskService) GetTasksByAvailableTaskType(userId uuid.UUID) ([]TaskDTO, error)
func (s *TaskService) GetTasksByCancelledTaskType(userId uuid.UUID) ([]TaskDTO, error)
func (s *TaskService) GetTasksByTaskType(userId uuid.UUID, taskType string) ([]TaskDTO, error)
func (s *TaskService) Create(repository.InsertNewTaskParams) (TaskDTO, error)
func (s *TaskService) Delete(taskId uuid.UUID) error
func (s *TaskService) UpdateAsAmbiguous(uuid.UUID) (TaskDTO, error)
func (s *TaskService) UpdateAsCancelled(uuid.UUID) (TaskDTO, error)
func (s *TaskService) UpdateAsDone(uuid.UUID) (TaskDTO, error)
func (s *TaskService) UpdateAsHold(uuid.UUID) (TaskDTO, error)
func (s *TaskService) UpdateAsPending(repository.UpdateAsPendingParams) (TaskDTO, error)
func (s *TaskService) UpdateAsRecurring(uuid.UUID) (TaskDTO, error)
func (s *TaskService) UpdateAsUndone(repository.UpdateAsUndoneParams) (TaskDTO, error)
func (s *TaskService) UpdateAsUrgent(userId uuid.UUID) (TaskDTO, error)
func (s *TaskService) UpdateTaskContent(repository.UpdateTaskContentParams) (TaskDTO, error)
func (s *TaskService) DeleteTaskContent(uuid.UUID) (TaskDTO, error)

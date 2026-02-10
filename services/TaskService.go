package services

import (
	"context"

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

func (s *TaskService) GetById(taskId uuid.UUID) (TaskDTO, error) {
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

// region TaskQueues
func (s *TaskService) GetTasksByScheduledTaskType(userId uuid.UUID) ([]TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return []TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	tasks, err := repo.GetTasksByScheduledTaskType(s.ctx, userId)
	if err != nil {
		return []TaskDTO{}, err
	}
	taskObjects, err := s.getTaskTypesForSlice(tasks, repo)
	tx.Commit(s.ctx)
	return taskObjects, nil
}

func (s *TaskService) GetTasksByCompletedTaskType(userId uuid.UUID) ([]TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return []TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	tasks, err := repo.GetTasksByCompletedTaskType(s.ctx, userId)
	if err != nil {
		return []TaskDTO{}, err
	}
	taskObjects, err := s.getTaskTypesForSlice(tasks, repo)
	tx.Commit(s.ctx)
	return taskObjects, nil
}

func (s *TaskService) GetTasksByAvailableTaskType(userId uuid.UUID) ([]TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return []TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	tasks, err := repo.GetTasksByAvailableTaskType(s.ctx, userId)
	if err != nil {
		return []TaskDTO{}, err
	}
	taskObjects, err := s.getTaskTypesForSlice(tasks, repo)
	tx.Commit(s.ctx)
	return taskObjects, nil
}

func (s *TaskService) GetTasksByCancelledTaskType(userId uuid.UUID) ([]TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return []TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	tasks, err := repo.GetTasksByCancelledTaskType(s.ctx, userId)
	if err != nil {
		return []TaskDTO{}, err
	}
	taskObjects, err := s.getTaskTypesForSlice(tasks, repo)
	tx.Commit(s.ctx)
	return taskObjects, nil
}

func (s *TaskService) GetTasksByTaskType(params repository.GetTasksByTaskTypeParams) ([]TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return []TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	tasks, err := repo.GetTasksByTaskType(s.ctx, params)
	if err != nil {
		return []TaskDTO{}, err
	}
	taskObjects, err := s.getTaskTypesForSlice(tasks, repo)
	tx.Commit(s.ctx)
	return taskObjects, nil
}
// endregion

// region Task Mutations
func (s *TaskService) Create(params repository.InsertNewTaskParams) (TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	task, err := repo.InsertNewTask(s.ctx, params)
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

func (s *TaskService) Delete(taskId uuid.UUID) error {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	err = repo.DeleteTask(s.ctx, taskId)
	if err != nil {
		return err
	}
	tx.Commit(s.ctx)
	return nil
}

func (s *TaskService) UpdateAsAmbiguous(taskId uuid.UUID) (TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	task, err := repo.UpdateAsAmbiguous(s.ctx, taskId)
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

func (s *TaskService) UpdateAsCancelled(taskId uuid.UUID) (TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	task, err := repo.UpdateAsCancelled(s.ctx, taskId)
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

func (s *TaskService) UpdateAsDone(taskId uuid.UUID) (TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	task, err := repo.UpdateAsDone(s.ctx,taskId)
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

func (s *TaskService) UpdateAsHold(taskId uuid.UUID) (TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	task, err := repo.UpdateAsHold(s.ctx,taskId)
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

func (s *TaskService) UpdateAsPending(params repository.UpdateAsPendingParams) (TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	task, err := repo.UpdateAsPending(s.ctx, params)
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
func (s *TaskService) UpdateAsRecurring(taskId uuid.UUID) (TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	task, err := repo.UpdateAsRecurring(s.ctx,taskId)
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

func (s *TaskService) UpdateAsUndone(params repository.UpdateAsUndoneParams) (TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	task, err := repo.UpdateAsUndone(s.ctx, params)
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
func (s *TaskService) UpdateAsUrgent(userId uuid.UUID) (TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	task, err := repo.UpdateAsUrgent(s.ctx,userId)
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

func (s *TaskService) UpdateTaskContent(params repository.UpdateTaskContentParams) (TaskDTO, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return TaskDTO{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	task, err := repo.UpdateTaskContent(s.ctx, params)
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
// endregion
